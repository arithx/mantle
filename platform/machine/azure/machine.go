// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in campliance with the License.
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
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/api/azure"
)

type machine struct {
	cluster *cluster
	mach    *azure.Machine
	dir     string
	journal *platform.Journal
	console []byte
}

func (am *machine) ID() string {
	return am.mach.ID
}

func (am *machine) IP() string {
	return am.mach.PublicIPAddress
}

func (am *machine) PrivateIP() string {
	return am.mach.PrivateIPAddress
}

func (am *machine) ResourceGroup() string {
	return am.cluster.ResourceGroup
}

func (am *machine) InterfaceName() string {
	return am.mach.InterfaceName
}

func (am *machine) PublicIPName() string {
	return am.mach.PublicIPName
}

func (am *machine) SSHClient() (*ssh.Client, error) {
	return am.cluster.SSHClient(am.IP())
}

func (am *machine) PasswordSSHClient(user string, password string) (*ssh.Client, error) {
	return am.cluster.PasswordSSHClient(am.IP(), user, password)
}

func (am *machine) SSH(cmd string) ([]byte, []byte, error) {
	return am.cluster.SSH(am, cmd)
}

func (am *machine) Reboot() error {
	err := platform.RebootMachine(am, am.journal, am.cluster.RuntimeConf())
	if err != nil {
		return err
	}

	// Re-fetch the Public & Private IP address for the event that it's changed during the reboot
	/*
	am.mach.PublicIPAddress, err = am.cluster.api.GetPublicIP(am.ID(), am.ResourceGroup())
	if err != nil {
		return err
	}
	am.mach.PrivateIPAddress, err = am.cluster.api.GetPrivateIP(am.InterfaceName(), am.ResourceGroup())
	*/
	am.mach.PublicIPAddress, am.mach.PrivateIPAddress, err = am.cluster.api.GetIPAddresses(am.InterfaceName(), am.PublicIPName(), am.ResourceGroup())
	return err
}

func (am *machine) Destroy() error {
	if err := am.saveConsole(); err != nil {
		// log error, but do not fail to terminate instance
		plog.Error(err)
	}

	if err := am.cluster.api.TerminateInstance(am.ID(), am.ResourceGroup()); err != nil {
		return fmt.Errorf("terminating instance: %v", err)
	}

	if am.journal != nil {
		if err := am.journal.Destroy(); err != nil {
			return err
		}
	}

	am.cluster.DelMach(am)

	return nil
}

func (am *machine) ConsoleOutput() string {
	return string(am.console)
}

func (am *machine) saveConsole() error {
	var err error
	am.console, err = am.cluster.api.GetConsoleOutput(am.ID(), am.ResourceGroup())
	if err != nil {
		return err
	}

	path := filepath.Join(am.dir, "console.txt")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(am.console)
	if err != nil {
		return fmt.Errorf("failed writing console to file: %v", err)
	}

	return nil
}
