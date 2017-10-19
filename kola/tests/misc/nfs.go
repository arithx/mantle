// Copyright 2015 CoreOS, Inc.
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

package misc

import (
	"fmt"
	"path"
	"time"

	"github.com/vincent-petithory/dataurl"

	"github.com/coreos/mantle/kola/cluster"
	"github.com/coreos/mantle/kola/register"
	"github.com/coreos/mantle/platform/conf"
	"github.com/coreos/mantle/util"
)

var (
	nfsserverconf = conf.Ignition(fmt.Sprintf(`{
		"ignition": {
			"version": "2.1.0"
		},
		"storage": {
			"files": [{
				"filesystem": "root",
				"path": "/etc/exports",
				"contents": { "source": "data:,%s" },
				"user": {"name": "core"},
				"group": {"name": "core"}
			}, {
				"filesystem": "root",
				"path": "/etc/hostname",
				"contents": { "source": "data:,nfs1" },
				"mode": 511
			}]
		},
		"systemd": {
			"units": [{
				"name": "rpc-statd.service",
				"enabled": true
			}, {
				"name": "rpc-mountd.service",
				"enabled": true
			}, {
				"name": "nfsd.service",
				"enabled": true
			}, {
				"name": "start-the-services.service",
				"enabled": true,
				"contents": "[Unit]\nAfter=rpc-statd.service\nRequires=rpc-statd.service\nAfter=rpc-mountd.service\nRequires=rpc-mountd.service\nAfter=nfsd.service\nRequires=nfsd.service\n\n[Service]\nExecStart=/usr/bin/echo start\n\n[Install]\nWantedBy=multi-user.target"
			}]
		}
	}`, dataurl.EscapeString("/tmp  *(ro,insecure,all_squash,no_subtree_check,fsid=0)")))

	mounttmpl = `[Unit]\nDescription=NFS Client\nAfter=network-online.target\nRequires=network-online.target\nAfter=rpc-statd.service\nRequires=rpc-statd.service\n\n[Mount]\nWhat=%s:/tmp\nWhere=/mnt\nType=nfs\nOptions=defaults,noexec,nfsvers=%d\n\n[Install]\nWantedBy=multi-user.target`
)

func init() {
	register.Register(&register.Test{
		Run:              NFSv3,
		ClusterSize:      0,
		Name:             "linux.nfs.v3",
	})
	register.Register(&register.Test{
		Run:              NFSv4,
		ClusterSize:      0,
		Name:             "linux.nfs.v4",
	})
}

func testNFS(c cluster.TestCluster, nfsversion int) {
	m1, err := c.NewMachine(nfsserverconf)
	if err != nil {
		c.Fatalf("Cluster.NewMachine: %s", err)
	}

	defer m1.Destroy()

	c.Log("NFS server booted.")

	/* poke a file in /tmp */
	tmp, err := c.SSH(m1, "mktemp")
	if err != nil {
		c.Fatalf("Machine.SSH: %s", err)
	}

	c.Logf("Test file %q created on server.", tmp)

	c2 := conf.Ignition(fmt.Sprintf(`{
		"ignition": {
			"version": "2.1.0"
		},
		"storage": {
			"files": [{
				"filesystem": "root",
                                "path": "/etc/hostname",
                                "contents": { "source": "data:,nfs2" },
				"mode": 511
			}]
		},
		"systemd": {
			"units": [{
				"name": "mnt.mount",
				"enabled": true,
				"contents": "%s"
			}]
		}
	}`, fmt.Sprintf(mounttmpl, m1.PrivateIP(), nfsversion)))

	m2, err := c.NewMachine(c2)
	if err != nil {
		c.Fatalf("Cluster.NewMachine: %s", err)
	}

	defer m2.Destroy()

	c.Log("NFS client booted.")

	checkmount := func() error {
		status, err := c.SSH(m2, "systemctl is-active mnt.mount")
		if err != nil || string(status) != "active" {
			return fmt.Errorf("mnt.mount status is %q: %v", status, err)
		}

		c.Log("Got NFS mount.")
		return nil
	}

	if err = util.Retry(10, 3*time.Second, checkmount); err != nil {
		c.Fatal(err)
	}

	_, err = c.SSH(m2, fmt.Sprintf("stat /mnt/%s", path.Base(string(tmp))))
	if err != nil {
		c.Fatalf("file %q does not exist", tmp)
	}
}

// Test that the kernel NFS server and client work within CoreOS.
func NFSv3(c cluster.TestCluster) {
	testNFS(c, 3)
}

// Test that NFSv4 without security works on CoreOS.
func NFSv4(c cluster.TestCluster) {
	testNFS(c, 4)
}
