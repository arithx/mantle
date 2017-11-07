// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/coreos/pkg/capnslog"
	"github.com/coreos/pkg/multierror"

	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/api/azure"
	"github.com/coreos/mantle/platform/conf"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "platform/machine/azure")
)

type cluster struct {
	*platform.BaseCluster
	api           *azure.API
	sshKey        string
	ResourceGroup string
}

func NewCluster(opts *azure.Options, rconf *platform.RuntimeConfig) (platform.Cluster, error) {
	api, err := azure.New(opts)
	if err != nil {
		return nil, err
	}

	bc, err := platform.NewBaseCluster(opts.BaseName, rconf, "")
	if err != nil {
		return nil, err
	}

	resourceGroup, err := api.CreateResourceGroup("kola-cluster")
	if err != nil {
		return nil, err
	}

	_, err = api.PrepareNetworkResources(resourceGroup)
	if err != nil {
		return nil, err
	}

	ac := &cluster{
		BaseCluster:   bc,
		api:           api,
		ResourceGroup: resourceGroup,
	}

	if !rconf.NoSSHKeyInMetadata {
		keys, err := ac.Keys()
		if err != nil {
			return nil, err
		}
		ac.sshKey = keys[0].String()
	} else {
		key, err := platform.GenerateFakeKey()
		if err != nil {
			return nil, err
		}
		ac.sshKey = key
	}

	return ac, nil
}

func (ac *cluster) vmname() string {
	b := make([]byte, 5)
	rand.Read(b)
	return fmt.Sprintf("%s-%x", ac.Name()[0:13], b)
}

func (ac *cluster) NewMachine(userdata *conf.UserData) (platform.Machine, error) {
	conf, err := ac.RenderUserData(userdata, map[string]string{
		"$private_ipv4": "${COREOS_AZURE_IPV4_DYNAMIC}",
		"$public_ipv4": "${COREOS_AZURE_IPV4_DYNAMIC}",
	})
	if err != nil {
		return nil, err
	}

	instance, err := ac.api.CreateInstance(ac.vmname(), conf.String(), ac.sshKey, ac.ResourceGroup)
	if err != nil {
		return nil, err
	}

	mach := &machine{
		cluster: ac,
		mach:    instance,
	}

	mach.dir = filepath.Join(ac.RuntimeConf().OutputDir, mach.ID())
	if err := os.Mkdir(mach.dir, 0777); err != nil {
		mach.Destroy()
		return nil, err
	}

	confPath := filepath.Join(mach.dir, "user-data")
	if err := conf.WriteFile(confPath); err != nil {
		mach.Destroy()
		return nil, err
	}

	if mach.journal, err = platform.NewJournal(mach.dir); err != nil {
		mach.Destroy()
		return nil, err
	}

	if err := platform.StartMachine(mach, mach.journal, ac.RuntimeConf()); err != nil {
		mach.Destroy()
		return nil, err
	}

	ac.AddMach(mach)

	return mach, nil
}

func (ac *cluster) Destroy() error {
	var err multierror.Error
	if e := ac.BaseCluster.Destroy(); e != nil {
		err = append(err, e)
	}
	if e := ac.api.TerminateResourceGroup(ac.ResourceGroup); e != nil {
		err = append(err, e)
	}
	return err.AsError()
}
