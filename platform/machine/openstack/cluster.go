// Copyright 2018 CoreOS, Inc.
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

package openstack

import (
	"os"
	"path/filepath"

	ctplatform "github.com/coreos/container-linux-config-transpiler/config/platform"
	"github.com/coreos/pkg/capnslog"

	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/api/openstack"
	"github.com/coreos/mantle/platform/conf"
)

const (
	Platform platform.Name = "openstack"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "platform/machine/openstack")
)

type cluster struct {
	*platform.BaseCluster
	api *openstack.API
}

// NewCluster creates an instance of a Cluster suitable for spawning
// instances on OpenStack.
func NewCluster(opts *openstack.Options, rconf *platform.RuntimeConfig) (platform.Cluster, error) {
	api, err := openstack.New(opts)
	if err != nil {
		return nil, err
	}

	bc, err := platform.NewBaseCluster(opts.Options, rconf, Platform, ctplatform.OpenStackMetadata)
	if err != nil {
		return nil, err
	}

	oc := &cluster{
		BaseCluster: bc,
		api:         api,
	}

	if !rconf.NoSSHKeyInMetadata {
		keys, err := oc.Keys()
		if err != nil {
			return nil, err
		}

		if err := api.AddKey(bc.Name(), keys[0].String()); err != nil {
			return nil, err
		}
	}

	return oc, nil
}

func (oc *cluster) NewMachine(userdata *conf.UserData) (platform.Machine, error) {
	conf, err := oc.RenderUserData(userdata, map[string]string{
		"$public_ipv4":  "${COREOS_OPENSTACK_IPV4_PUBLIC}",
		"$private_ipv4": "${COREOS_OPENSTACK_IPV4_LOCAL}",
	})
	if err != nil {
		return nil, err
	}

	var keyname string
	if !oc.RuntimeConf().NoSSHKeyInMetadata {
		keyname = oc.Name()
	}
	instance, err := oc.api.CreateInstance(oc.Name(), keyname, conf.String())
	if err != nil {
		return nil, err
	}

	mach := &machine{
		cluster: oc,
		mach:    instance,
	}

	mach.dir = filepath.Join(oc.RuntimeConf().OutputDir, mach.ID())
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

	if err := platform.StartMachine(mach, mach.journal); err != nil {
		mach.Destroy()
		return nil, err
	}

	oc.AddMach(mach)

	return mach, nil
}

func (oc *cluster) Destroy() {
	if !oc.RuntimeConf().NoSSHKeyInMetadata {
		if err := oc.api.DeleteKey(oc.Name()); err != nil {
			plog.Errorf("Error deleting key %v: %v", oc.Name(), err)
		}
	}

	oc.BaseCluster.Destroy()
}
