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

package update

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/coreos/go-omaha/omaha"

	"github.com/coreos/mantle/kola"
	"github.com/coreos/mantle/kola/cluster"
	"github.com/coreos/mantle/kola/register"
	tutil "github.com/coreos/mantle/kola/tests/util"
	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/conf"
	"github.com/coreos/mantle/platform/local"
	"github.com/coreos/mantle/util"
)

type offsetValue struct {
	offset  int
	value   []byte
	newHash string
}

var (
	grubUpdates = map[string]offsetValue{
		"1e11052c144ae483cba4f70efe278070c50daa80d1a5febe7c0d08e401baf0ab16a542b7a34da00d7aea1591238f01d9902ac6fbfcce8d82eebcf09d97d132cd": {3378, []byte{0x74, 2, 0, 0},
			"3e9cc58b60301637e37e504c6555f5758c79e3c5493478601ca3c77cc4b6186aa9d732de268a6feb7fcbb88190b46d186b3e3d2dec2ac418645dc1b86a026c0d"},
		"2fd5f0fade7c4c986259524f148f79ee1d1353d7ab83d1bdd0d50e52d393d8d896c32ab64eb714ba08861b8ba4f113d19f940a04889fa407784f010f119c8c19": {3378, []byte{0x88, 2, 0, 0},
			"5ce27d73800d4075e4fc48928d7caa11d070b50c4ff056fc56fa93102878cecafa68513a22b71c8e68c749f277d9004bb9e0b1d3394cdb78ea2383b0658e4a28"},
		"3c591be75a6aa903ee8deed5f8116e627f53738eb8a2ceb80aeaab08f485b405b87565148e482a7234138302eb900786d0e14939a3c3451d424052ca2bd73181": {3383, []byte{0x74, 2, 0, 0},
			"0a385b66953fe4125ba3e8761a718ce3874c7d63363c0f74c09f897af8971bdd67633921c11bf65d2c7e19346b176e6493726e58ad54b64d4c336529c6a34497"},
		"6d60e369c1b4b484c7221e91d80f03a782b5286137f8087b2bf22b9f54e3507c4947d2a456ba46a6ef3c0c7216dc8251c017e9122f44cba89e32f23a0542afd3": {3383, []byte{0x74, 2, 0, 0},
			"a70eb783b1f9d1b8d97dcee3f98787a57a9f49ceffa02abe7ab3ca296824f2f8b0db8f8646130a328ad5e2f11d5dd061e5a0d85d8390f722b5d573c62a1307b0"},
		"6e9e5ebb6cd1a15d5a570d9a06a56c9bf60cb047d858c7220dd9dcfd54ebb87c8e0ea4611ac98a0ab51fec9e87b265b2207973f3e4882d87f62d887da72f87ab": {3404, []byte{0x88, 2, 0, 0},
			"937e579767de1f9a7869fbe9b25e296349c4ad51dda30bfd98932e450ee52854a9019eeb75dea2a1f0eff6ba9f06bec794ff05c9e738c2ffd7d5f0eaa751eb9c"},
		"82b37fa4b305cab33277d2cf0249008731a69575b5689a47e72fe2a35be4440e0e116bc02191f9b0066ea3ae278327fe3409f28d25d13bae88c5f347dba6a254": {3383, []byte{0x74, 2, 0, 0},
			"a272cf4995e6406f3652be38aa8b971b928427bbe855f266585ddfe4e9c6e79c926a1665d1ec26ac83a0cefafc9f7c0fadbb433ed8748493df14b37174dde37b"},
		"8a7b03d92a8b115943e7f004820fadd2dc6ab125c077a48fb232a1e9ac77fdb27fbb01d52fd33a6ddf65a9f58ce981244c99bcca821030511caa277bc2f68239": {3378, []byte{0x74, 2, 0, 0},
			"e61536f2e47d95cc98cb37e955c8cb4b1afaa67f4705868e9f0f5f657370d2a01ea7cfe2733e0d7f0295e18282e24090f82fbbfa88abf6400f53ba476a567bdd"},
		"a3e9dadfe3cc34189b5fee83bfc01c3c5b42e04ba19cdcd84f8301c42566617b5916294e2414348139d8c5e557a7ccf6c0d3dca0661f2d10c0c0077345630b1d": {3298, []byte{0x74, 2, 0, 0},
			"c9be6791c064167d2b43ac71c2035a7501542901e0f1a4050466b016d2b85ea46da5780e5b9b61bb26e44ccb39e58fa461d1d9cc3ffb5be26aa35458adff542b"},
		"c127d7c1dbd5d11cf7af627e37808ea16166b6430ddd8e96111e503cc78ae1fd78083d474495951743fa1b489140be63178a4bb65dabb0d719c5d0ad9c57eb78": {3298, []byte{0x74, 2, 0, 0},
			"d7c74e09f7030ebcdcb0c4b5e99cb33fbab1138dd552ee194eb79f71cf5252ef4e169a366a160301fcf6a4e8265912c4b334019e47b74c7b2126a605a910eb6f"},
		"f1f9abefa49eeba6a3fe46ba3d254dfc3fa6e2cd8823835e2d982e9cbcd0d82c298e3896dc79d234305d2262a828a7398526906022f8ed7407368725d95e08d8": {3375, []byte{0x74, 2, 0, 0},
			"daf12ef71b17d7ddc8fb8419b19bbe2ffb61be9b5bd0d1ca320ac8a07fbfe49924a4b1167bd356d96218972b90b936ce0e508d3b08d470945866aa0c7e71510d"},
	}
)

func init() {
	register.Register(&register.Test{
		Name:        "coreos.update.payload",
		Run:         payload,
		ClusterSize: 1,
		NativeFuncs: map[string]func() error{
			"Omaha": Serve,
		},
	})
	register.Register(&register.Test{
		Name:        "coreos.update.payload.grub",
		Run:         grub,
		ClusterSize: 1,
		UserData:    conf.CloudConfig(""),
		NativeFuncs: map[string]func() error{
			"Omaha": Serve,
		},
	})
}

func Serve() error {
	omahaserver, err := omaha.NewTrivialServer(":34567")
	if err != nil {
		return fmt.Errorf("creating trivial omaha server: %v\n", err)
	}

	omahawrapper := local.OmahaWrapper{TrivialServer: omahaserver}

	if err = omahawrapper.AddPackage("/updates/update.gz", "update.gz"); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}

	return omahawrapper.Serve()
}

func gunzipAndRead(comp []byte) ([]byte, error) {
	fh := bytes.NewReader(comp)

	uncomp, err := gzip.NewReader(fh)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(uncomp)
}

func grub(c cluster.TestCluster) {
	addr := configureOmahaServer(c, c.Machines()[0])

	// create the actual test machine, the machine
	// that is created by the test registration is
	// used to host the omaha server
	m, err := c.NewMachine(conf.CloudConfig(""))
	if err != nil {
		c.Fatalf("creating test machine: %v", err)
	}

	// Machines are intentionally configured post-boot
	// via SSH to allow for testing versions which predate
	// Ignition
	configureMachineForUpdate(c, m, addr)

	originalBytes, err := gunzipAndRead(c.MustSSH(m, "cat /boot/coreos/grub/i386-pc/linux.mod"))
	if err != nil {
		c.Fatalf("failed decompressing: %v", err)
	}
	sumAr := sha512.Sum512(originalBytes)
	sum := hex.EncodeToString(sumAr[:])
	offsetValue, ok := grubUpdates[sum]
	if !ok {
		c.Fatalf("did not find bad linux.mod")
	}

	validateGrub := func() {
		newBytes, err := gunzipAndRead(c.MustSSH(m, "cat /boot/coreos/grub/i386-pc/linux.mod"))
		if err != nil {
			c.Fatalf("failed decompressing: %v", err)
		}
		// golang doesn't do slice comparison. golang --
		if !reflect.DeepEqual(newBytes[offsetValue.offset:offsetValue.offset+4], offsetValue.value) {
			c.Fatalf("did not patch grub")
		}
		sumArInner := sha512.Sum512(newBytes)
		if offsetValue.newHash != hex.EncodeToString(sumArInner[:]) {
			c.Fatalf("final hash did not match: expected %v got %v", offsetValue.newHash, hex.EncodeToString(sumArInner[:]))
		}

		c.MustSSH(m, "cat /boot/coreos/grub/skip-bug-2400-patch")
	}

	updateMachine(c, m)

	tutil.AssertBootedUsr(c, m, "USR-B")

	validateGrub()

	tutil.InvalidateUsrPartition(c, m, "USR-A")

	updateMachine(c, m)

	tutil.AssertBootedUsr(c, m, "USR-A")

	validateGrub()
}

func payload(c cluster.TestCluster) {
	addr := configureOmahaServer(c, c.Machines()[0])

	// create the actual test machine, the machine
	// that is created by the test registration is
	// used to host the omaha server
	m, err := c.NewMachine(nil)
	if err != nil {
		c.Fatalf("creating test machine: %v", err)
	}

	// Machines are intentionally configured post-boot
	// via SSH to allow for testing versions which predate
	// Ignition
	configureMachineForUpdate(c, m, addr)

	tutil.AssertBootedUsr(c, m, "USR-A")

	updateMachine(c, m)

	tutil.AssertBootedUsr(c, m, "USR-B")

	tutil.InvalidateUsrPartition(c, m, "USR-A")

	updateMachine(c, m)

	tutil.AssertBootedUsr(c, m, "USR-A")
}

func configureOmahaServer(c cluster.TestCluster, srv platform.Machine) string {
	if kola.UpdatePayloadFile == "" {
		c.Skip("no update payload provided")
	}

	in, err := os.Open(kola.UpdatePayloadFile)
	if err != nil {
		c.Fatalf("opening update payload: %v", err)
	}
	defer in.Close()
	if err := platform.InstallFile(in, srv, "/updates/update.gz"); err != nil {
		c.Fatalf("copying update payload to omaha server: %v", err)
	}

	platform.Manhole(srv)
	c.MustSSH(srv, fmt.Sprintf("sudo systemd-run --quiet ./kolet run %s Omaha", c.Name()))

	err = util.WaitUntilReady(60*time.Second, 5*time.Second, func() (bool, error) {
		_, _, err := srv.SSH(fmt.Sprintf("curl %s:34567", srv.PrivateIP()))
		return err == nil, nil
	})
	if err != nil {
		c.Fatal("timed out waiting for omaha server to become active")
	}

	return fmt.Sprintf("%s:34567", srv.PrivateIP())
}

func configureMachineForUpdate(c cluster.TestCluster, m platform.Machine, addr string) {
	// update atomicly so nothing reading update.conf fails
	c.MustSSH(m, fmt.Sprintf(`sudo bash -c "cat >/etc/coreos/update.conf.new <<EOF
GROUP=developer
SERVER=http://%s/v1/update
EOF"`, addr))
	c.MustSSH(m, "sudo mv /etc/coreos/update.conf{.new,}")

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

	// disable reboot so the test has explicit control
	c.MustSSH(m, "sudo systemctl mask --now locksmithd.service")
	c.MustSSH(m, "sudo systemctl reset-failed locksmithd.service")

	c.MustSSH(m, "sudo systemctl restart update-engine.service")
}

func updateMachine(c cluster.TestCluster, m platform.Machine) {
	c.Logf("Triggering update_engine")

	out, stderr, err := m.SSH("update_engine_client -check_for_update")
	if err != nil {
		c.Fatalf("Executing update_engine_client failed: %v: %v: %s", out, err, stderr)
	}

	err = util.WaitUntilReady(240*time.Second, 10*time.Second, func() (bool, error) {
		envs, stderr, err := m.SSH("update_engine_client -status 2>/dev/null")
		if err != nil {
			return false, fmt.Errorf("checking status failed: %v: %s", err, stderr)
		}

		return splitNewlineEnv(string(envs))["CURRENT_OP"] == "UPDATE_STATUS_UPDATED_NEED_REBOOT", nil
	})
	if err != nil {
		c.Fatalf("waiting for UPDATE_STATUS_UPDATED_NEED_REBOOT: %v", err)
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
