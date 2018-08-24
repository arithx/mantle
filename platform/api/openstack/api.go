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
	"fmt"
	"strings"
	"time"

	"github.com/coreos/pkg/capnslog"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	//"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"

	"github.com/coreos/mantle/auth"
	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/util"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "platform/api/openstack")
)

type Options struct {
	*platform.Options

	// Config file. Defaults to $HOME/.config/openstack.json.
	ConfigPath string
	// Profile name
	Profile string

	// Region (e.g. "regionOne")
	Region string
	// Instance Flavor ID
	Flavor string
	// Image ID
	Image string
	// Network ID
	Network string
}

type Machine struct {
	Server     *servers.Server
	FloatingIP *floatingips.FloatingIP
}

type API struct {
	opts          *Options
	computeClient *gophercloud.ServiceClient
	//imageClient   *gophercloud.ServiceClient
	networkClient *gophercloud.ServiceClient
}

func New(opts *Options) (*API, error) {
	profiles, err := auth.ReadOpenStackConfig(opts.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read OpenStack config: %v", err)
	}

	if opts.Profile == "" {
		opts.Profile = "default"
	}
	profile, ok := profiles[opts.Profile]
	if !ok {
		return nil, fmt.Errorf("no such profile %q", opts.Profile)
	}

	osOpts := gophercloud.AuthOptions{
		IdentityEndpoint: profile.AuthURL,
		TenantID:         profile.TenantID,
		TenantName:       profile.TenantName,
		Username:         profile.Username,
		Password:         profile.Password,
	}

	provider, err := openstack.AuthenticatedClient(osOpts)
	if err != nil {
		return nil, fmt.Errorf("failed creating provider: %v", err)
	}

	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Name:   "nova",
		Region: opts.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %v", err)
	}

	/*
	imageClient, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Name:   "glance",
		Region: opts.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create image client: %v", err)
	}
	*/

	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: opts.Region,
	})

	a := &API{
		opts:          opts,
		computeClient: computeClient,
		//imageClient:   imageClient,
		networkClient: networkClient,
	}

	err = a.resolveFlavor()
	if err != nil {
		return nil, fmt.Errorf("resolving flavor: %v", err)
	}

	err = a.resolveImage()
	if err != nil {
		return nil, fmt.Errorf("resolving image: %v", err)
	}

	err = a.resolveNetwork()
	if err != nil {
		return nil, fmt.Errorf("resolving network: %v", err)
	}

	return a, nil
}

func (a *API) resolveFlavor() error {
	pager := flavors.ListDetail(a.computeClient, flavors.ListOpts{})
	if pager.Err != nil {
		return fmt.Errorf("retrieving flavors: %v", pager.Err)
	}

	pages, err := pager.AllPages()
	if err != nil {
		return fmt.Errorf("retrieving flavors pages: %v", err)
	}

	empty, err := pages.IsEmpty()
	if err != nil {
		return fmt.Errorf("parsing flavor pages: %v", err)
	}
	if empty {
		return fmt.Errorf("no flavors found")
	}

	flavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmt.Errorf("extracting flavors: %v", err)
	}

	for _, flavor := range flavors {
		if flavor.ID == a.opts.Flavor || flavor.Name == a.opts.Flavor {
			a.opts.Flavor = flavor.ID
			return nil
		}
	}

	return fmt.Errorf("specified flavor %q not found", a.opts.Flavor)
}

func (a *API) resolveImage() error {
	pager := images.ListDetail(a.computeClient, images.ListOpts{})
	if pager.Err != nil {
		return fmt.Errorf("retrieving images: %v", pager.Err)
	}

	pages, err := pager.AllPages()
	if err != nil {
		return fmt.Errorf("retrieving image pages: %v", err)
	}

	empty, err := pages.IsEmpty()
	if err != nil {
		return fmt.Errorf("parsing image pages: %v", err)
	}
	if empty {
		return fmt.Errorf("no images found")
	}

	images, err := images.ExtractImages(pages)
	if err != nil {
		return fmt.Errorf("extracting images: %v", err)
	}

	for _, image := range images {
		if image.ID == a.opts.Image || image.Name == a.opts.Image {
			a.opts.Image = image.ID
			return nil
		}
	}

	return fmt.Errorf("specified image %q not found", a.opts.Image)
}

func (a *API) resolveNetwork() error {
	networks, err := a.GetNetworks()
	if err != nil {
		return err
	}

	for _, network := range networks {
		if network.ID == a.opts.Network || network.Name == a.opts.Network {
			a.opts.Network = network.ID
			return nil
		}
	}

	return fmt.Errorf("specified network %q not found", a.opts.Network)
}

func (a *API) PreflightCheck() error {
	if err := servers.List(a.computeClient, servers.ListOpts{}).Err; err != nil {
		return fmt.Errorf("listing servers: %v", err)
	}
	return nil
}

func (a *API) CreateInstance(name, sshKeyID, userdata string) (machine *Machine, err error) {
	var networkID string
	if a.opts.Network == "" {
		networkID, err = a.GetFirstNetworkID()
		if err != nil {
			return nil, fmt.Errorf("getting network: %v", err)
		}
	} else {
		networkID = a.opts.Network
	}

	floatingip, err := a.CreateFloatingIP()
	if err != nil {
		return nil, fmt.Errorf("creating floating ip: %v", err)
	}

	var server *servers.Server
	server, err = servers.Create(a.computeClient, keypairs.CreateOptsExt{
		servers.CreateOpts{
			Name: name,
			FlavorRef: a.opts.Flavor,
			ImageRef: a.opts.Image,
			Metadata: map[string]string{
				"CreatedBy":    "mantle",
			},
			SecurityGroups: []string{"default", "wideopen"},
			Networks: []servers.Network{
				{
					UUID: networkID,
				},
			},
			UserData: []byte(userdata),
		},
		sshKeyID,
	}).Extract()
	if err != nil {
		a.DeleteFloatingIP(floatingip.ID)
		return nil, fmt.Errorf("creating server: %v", err)
	}

	serverID := server.ID
	
	err = util.WaitUntilReady(5*time.Minute, 10*time.Second, func() (bool, error) {
		var err error
		server, err = servers.Get(a.computeClient, serverID).Extract()
		if err != nil {
			return false, err
		}
		return server.Status == "ACTIVE", nil
	})
	if err != nil {
		a.DeleteInstance(serverID)
		a.DeleteFloatingIP(floatingip.ID)
		return nil, fmt.Errorf("waiting for instance to run: %v", err)
	}

	err = floatingips.AssociateInstance(a.computeClient, serverID, floatingips.AssociateOpts{
		FloatingIP: floatingip.IP,
	}).ExtractErr()
	if err != nil {
		a.DeleteInstance(serverID)
		a.DeleteFloatingIP(floatingip.ID)
		return nil, fmt.Errorf("associating floating ip: %v", err)
	}

	server, err = servers.Get(a.computeClient, serverID).Extract()
	if err != nil {
		a.DisassociateFloatingIP(serverID, floatingip.IP)
		a.DeleteInstance(serverID)
		a.DeleteFloatingIP(floatingip.ID)
		return nil, fmt.Errorf("retrieving server info: %v", err)
	}

	return &Machine{
		Server: server,
		FloatingIP: floatingip,
	}, nil
}

func (a *API) GetFirstNetworkID() (string, error) {
	networks, err := a.GetNetworks()
	if err != nil {
		return "", err
	}

	return networks[0].ID, nil
}

func (a *API) GetNetworks() ([]networks.Network, error) {
	pager := networks.List(a.networkClient, networks.ListOpts{})
	if pager.Err != nil {
		return nil, fmt.Errorf("retrieving networks: %v", pager.Err)
	}

	pages, err := pager.AllPages()
	if err != nil {
		return nil, fmt.Errorf("retrieving network pages: %v", err)
	}

	empty, err := pages.IsEmpty()
	if err != nil {
		return nil, fmt.Errorf("parsing network pages: %v", err)
	}
	if empty {
		return nil, fmt.Errorf("no networks found")
	}

	networks, err := networks.ExtractNetworks(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting networks: %v", err)
	}
	return networks, nil
}

func (a *API) CreateFloatingIP() (*floatingips.FloatingIP, error) {
	return floatingips.Create(a.computeClient, floatingips.CreateOpts{
		Pool: "10.8.240.0",
	}).Extract()
}

func (a *API) DisassociateFloatingIP(serverID, id string) error {
	return floatingips.DisassociateInstance(a.computeClient, serverID, floatingips.DisassociateOpts{
		FloatingIP: id,
	}).ExtractErr()
}

func (a *API) DeleteFloatingIP(id string) error {
	return floatingips.Delete(a.computeClient, id).ExtractErr()
}

func (a *API) DeleteInstance(id string) (err error) {
	return servers.Delete(a.computeClient, id).ExtractErr()
}

func (a *API) listInstancesWithMetadata(metadata map[string]string) ([]servers.Server, error) {
	pager := servers.List(a.computeClient, servers.ListOpts{})
	if pager.Err != nil {
		return nil, fmt.Errorf("retrieving servers: %v", pager.Err)
	}
	
	pages, err := pager.AllPages()
	if err != nil {
		return nil, fmt.Errorf("retrieving server pages: %v", err)
	}

	empty, err := pages.IsEmpty()
	if err != nil {
		return nil, fmt.Errorf("parsing server pages: %v", err)
	}
	if empty {
		return nil, nil
	}

	allServers, err := servers.ExtractServers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting servers: %v", err)
	}
	var retServers []servers.Server
	for _, server := range allServers {
		isMatch := true
		for key, val := range metadata {
			if value, ok := server.Metadata[key]; !ok || val != value {
				isMatch = false
			}
		}
		if isMatch {
			retServers = append(retServers, server)
		}
	}
	return retServers, nil
}

func (a *API) GetConsoleOutput(id string) (string, error) {
	return servers.ShowConsoleOutput(a.computeClient, id, servers.ShowConsoleOutputOpts{}).Extract()
}

/*
func (a *API) DeleteImage(imageID string) error {
	return imageservice.Delete(a.imageClient, imageID).ExtractErr()
}
*/

func (a *API) AddKey(name, key string) error {
	_, err := keypairs.Create(a.computeClient, keypairs.CreateOpts{
		Name:      name,
		PublicKey: key,
	}).Extract()
	return err
}

func (a *API) DeleteKey(name string) error {
	return keypairs.Delete(a.computeClient, name).ExtractErr()
}


func (a *API) GC(gracePeriod time.Duration) error {
	threshold := time.Now().Add(-gracePeriod)

	servers, err := a.listInstancesWithMetadata(map[string]string{
		"CreatedBy": "mantle",
	})
	if err != nil {
		return err
	}
	for _, server := range servers {
		if !strings.Contains(server.Status, "DELETED") {
			if server.Created.After(threshold) {
				continue
			}
			if err := a.DeleteInstance(server.ID); err != nil {
				return fmt.Errorf("couldn't delete server %s: %v", server.ID, err)
			}
		}
	}
	return nil
}
