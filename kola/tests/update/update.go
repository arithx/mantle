// Copyright 2016 CoreOS, Inc.
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

package update

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-omaha/omaha"

	"github.com/coreos/mantle/kola"
	"github.com/coreos/mantle/kola/cluster"
	"github.com/coreos/mantle/kola/register"
	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/conf"
	"github.com/coreos/mantle/platform/local"
	"github.com/coreos/mantle/platform/machine/qemu"
	"github.com/coreos/mantle/sdk"
	"github.com/coreos/mantle/util"
)

func init() {
	register.Register(&register.Test{
		Name:        "coreos.update.updatepayload",
		Run:         updatepayload,
		ClusterSize: 1,
		NativeFuncs: map[string]func() error {
			"Omaha": Serve,
		},
		Platforms:   []string{"qemu", "aws"},
	})
}

func Serve() error {
	go func() {
		omahaserver, err := omaha.NewTrivialServer(":34567")
		if err != nil {
			fmt.Printf("creating trivial omaha server: %v\n", err)
		}

		omahawrapper := local.OmahaWrapper{TrivialServer: omahaserver}

		if err = omahawrapper.AddPackage("/updates/update.gz", "update.gz"); err != nil {
			fmt.Printf("bad payload: %v", err)
		}

		omahawrapper.TrivialServer.SetVersion("9999.9.9")

		go omahawrapper.Serve()
	}()

	select {}
}

func updatepayload(c cluster.TestCluster) {
	// create the actual test machine, the machine
	// that is created by the test registration is
	// used to host the omaha server if running on
	// a non-qemu platform
	m, err := c.NewMachine(conf.Ignition(`{"ignition":{"version": "2.1.0"}}`))
	if err != nil {
		c.Fatalf("creating test machine: %v", err)
	}

	addr := configureOmahaServer(c)

	configureMachine(m, c, addr)

	checkUsrA(m, c)

	updateMachine(m, c)

	checkUsrB(m, c)

	if out, stderr, err := m.SSH("sudo coreos-setgoodroot && sudo wipefs /dev/disk/by-partlabel/USR-A"); err != nil {
		c.Fatalf("invalidating USR-A failed: %s: %v: %s", out, err, stderr)
	}

	updateMachine(m, c)

	checkUsrA(m, c)
}

func configureOmahaServer(c cluster.TestCluster) string {
	if qc, ok := c.Cluster.(*qemu.Cluster); ok {
		if err := qc.OmahaServer.AddPackage(kola.UpdatePayload, "update.gz"); err != nil {
			c.Fatalf("bad payload: %v", err)
		}
		qc.OmahaServer.SetVersion("9999.9.9")

		port, ok := qc.OmahaServer.Addr().(*net.TCPAddr)
		if !ok {
			c.Fatal("failed detecting update server port")
		}

		return fmt.Sprintf("10.0.0.1:%d", port.Port)
	} else {
		srv := c.Machines()[0]

		in, err := os.Open(kola.UpdatePayload)
		if err != nil {
			c.Fatalf("opening update payload: %v", err)
		}
		defer in.Close()
		if err := platform.InstallFile(in, srv, "/updates/update.gz"); err != nil {
			c.Fatalf("copying update payload to omaha server: %v", err)
		}

		c.MustSSH(srv, fmt.Sprintf("sudo systemd-run --quiet ./kolet run %s Omaha", c.Name()))

		return fmt.Sprintf("%s:34567", srv.PrivateIP())
	}
}

func configureMachine(m platform.Machine, c cluster.TestCluster, addr string) {
	// update atomicly so nothing reading update.conf fails
	c.MustSSH(m, fmt.Sprintf(`sudo bash -c "cat >/etc/coreos/update.conf.new <<EOF
GROUP=developer
SERVER=http://%s/v1/update
EOF"`, addr))
	c.MustSSH(m, "sudo mv /etc/coreos/update.conf{.new,}")

	/*
	// inject dev key
	c.MustSSH(m, `sudo bash -c "cat >/etc/coreos/update-payload-key.pub.pem <<EOF
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzFS5uVJ+pgibcFLD3kbY
k02Edj0HXq31ZT/Bva1sLp3Ysv+QTv/ezjf0gGFfASdgpz6G+zTipS9AIrQr0yFR
+tdp1ZsHLGxVwvUoXFftdapqlyj8uQcWjjbN7qJsZu0Ett/qo93hQ5nHW7Sv5dRm
/ZsDFqk2Uvyaoef4bF9r03wYpZq7K3oALZ2smETv+A5600mj1Xg5M52QFU67UHls
EFkZphrGjiqiCdp9AAbAvE7a5rFcJf86YR73QX08K8BX7OMzkn3DsqdnWvLB3l3W
6kvIuP+75SrMNeYAcU8PI1+bzLcAG3VN3jA78zeKALgynUNH50mxuiiU3DO4DZ+p
5QIDAQAB
-----END PUBLIC KEY-----
EOF"`)

	c.MustSSH(m, "sudo mount --bind /etc/coreos/update-payload-key.pub.pem /usr/share/update_engine/update-payload-key.pub.pem")
	*/

	// disable reboot so the test has explicit control
	c.MustSSH(m, "sudo systemctl mask locksmithd.service")
	c.MustSSH(m, "sudo systemctl stop locksmithd.service")
	c.MustSSH(m, "sudo systemctl reset-failed locksmithd.service")

	c.MustSSH(m, "sudo systemctl restart update-engine.service")

	c.MustSSH(m, "update_engine --v=10")
}

func updateMachine(m platform.Machine, c cluster.TestCluster) {
	c.Logf("Triggering update_engine")

	out, stderr, err := m.SSH("update_engine_client -check_for_update")
	if err != nil {
		c.Fatalf("Executing update_engine_client failed: %v: %v: %s", out, err, stderr)
	}

	err = util.WaitUntilReady(120 * time.Second, 10 * time.Second, func() (bool, error) {
		envs, stderr, err := m.SSH("update_engine_client -status 2>/dev/null")
		if err != nil {
			return false, fmt.Errorf("checking status failed: %v: %s", err, stderr)
		}

		return splitNewlineEnv(string(envs))["CURRENT_OP"] == "UPDATE_STATUS_UPDATED_NEED_REBOOT", nil
	})
	if err != nil {
		c.Fatalf("Updating machine: %v", err)
	}

	c.Logf("Rebooting test machine")

	if err = m.Reboot(); err != nil {
		c.Fatalf("reboot failed: %v", err)
	}
}

// splits newline-delimited KEY=VAL pairs into a map
func splitNewlineEnv(envs string) map[string]string {
    m := make(map[string]string)
    sc := bufio.NewScanner(strings.NewReader(envs))
    for sc.Scan() {
        spl := strings.SplitN(sc.Text(), "=", 2)
        m[spl[0]] = spl[1]
    }
    return m
}

// split space-seperated KEY=VAL pairs into a map
func splitSpaceEnv(envs string) map[string]string {
    m := make(map[string]string)
    pairs := strings.Fields(envs)
    for _, p := range pairs {
        spl := strings.SplitN(p, "=", 2)
        if len(spl) == 2 {
            m[spl[0]] = spl[1]
        }
    }
    return m
}

func checkUsrA(m platform.Machine, c cluster.TestCluster) {
	c.Logf("Checking for boot from USR-A partition")
	checkUsrPartition(m, c, []string{
		"PARTUUID=" + sdk.USRAUUID.String(),
		"PARTLABEL=USR-A"})
}

func checkUsrB(m platform.Machine, c cluster.TestCluster) {
	c.Logf("Checking for boot from USR-B partition")
	checkUsrPartition(m, c, []string{
		"PARTUUID=" + sdk.USRBUUID.String(),
		"PARTLABEL=USR-B"})
}

func checkUsrPartition(m platform.Machine, c cluster.TestCluster, accept []string) {
	out, stderr, err := m.SSH("cat /proc/cmdline")
	if err != nil {
		c.Fatalf("cat /proc/cmdline: %v: %v: %s", out, err, stderr)
	}
	c.Logf("Kernel cmdline: %s", out)

	vars := splitSpaceEnv(string(out))
	for _, a := range accept {
		if vars["mount.usr"] == a {
			return
		}
		if vars["verity.usr"] == a {
			return
		}
		if vars["usr"] == a {
			return
		}
	}

	c.Fatalf("mount.usr not one of %q", strings.Join(accept, " "))
}
