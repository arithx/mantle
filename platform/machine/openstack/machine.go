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

	"golang.org/x/crypto/ssh"

	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/api/openstack"
)

type machine struct {
	cluster *cluster
	mach    *openstack.Machine
	dir     string
	journal *platform.Journal
	console string
}

func (om *machine) ID() string {
	return om.mach.Server.ID
}

func (om *machine) IP() string {
	if om.mach.FloatingIP != nil {
		return om.mach.FloatingIP.IP
	} else {
		return om.mach.Server.AccessIPv4
	}
}

func (om *machine) PrivateIP() string {
	return om.IP()
}

func (om *machine) RuntimeConf() platform.RuntimeConfig {
	return om.cluster.RuntimeConf()
}

func (om *machine) SSHClient() (*ssh.Client, error) {
	return om.cluster.SSHClient(om.IP())
}

func (om *machine) PasswordSSHClient(user string, password string) (*ssh.Client, error) {
	return om.cluster.PasswordSSHClient(om.IP(), user, password)
}

func (om *machine) SSH(cmd string) ([]byte, []byte, error) {
	return om.cluster.SSH(om, cmd)
}

func (om *machine) Reboot() error {
	return platform.RebootMachine(om, om.journal)
}

func (om *machine) Destroy() {
	if err := om.saveConsole(); err != nil {
		plog.Errorf("Error saving console for instance %v: %v", om.ID(), err)
	}

	if err := om.cluster.api.DisassociateFloatingIP(om.ID(), om.IP()); err != nil {
		plog.Errorf("Disassociating FloatingIP for instance %v: %v", om.ID(), err)
	}

	if err := om.cluster.api.DeleteInstance(om.ID()); err != nil {
		plog.Errorf("Error terminating instance %v: %v", om.ID(), err)
	}

	if om.mach.FloatingIP != nil {
		if err := om.cluster.api.DeleteFloatingIP(om.mach.FloatingIP.ID); err != nil {
			plog.Errorf("Error terminating FloatingIP %v: %v", om.mach.FloatingIP.ID, err)
		}
	}

	if om.journal != nil {
		om.journal.Destroy()
	}

	om.cluster.DelMach(om)
}

func (om *machine) ConsoleOutput() string {
	return om.console
}

func (om *machine) saveConsole() error {
	var err error
	om.console, err = om.cluster.api.GetConsoleOutput(om.ID())
	if err != nil {
		plog.Warningf("Error retrieving console log for %v: %v", om.ID(), err)
	}

	path := filepath.Join(om.dir, "console.txt")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(om.console)

	return nil
}

func (om *machine) JournalOutput() string {
	if om.journal == nil {
		return ""
	}

	data, err := om.journal.Read()
	if err != nil {
		plog.Errorf("Reading journal for instance %v: %v", om.ID(), err)
	}
	return string(data)
}
