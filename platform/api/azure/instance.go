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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/network"

	"github.com/coreos/mantle/util"
)

type Machine struct {
	ID string
	PublicIPAddress string
	PrivateIPAddress string
	InterfaceName string
	PublicIPName string
}

func (a *API) getVirtualMachine(name, resourceGroup string) (compute.VirtualMachine, error) {
	auth, err := auth.GetClientSetup(compute.DefaultBaseURI)
	if err != nil {
		return compute.VirtualMachine{}, err
	}
	client := compute.NewVirtualMachinesClientWithBaseURI(auth.BaseURI, auth.SubscriptionID)
	client.Authorizer = auth

	return client.Get(resourceGroup, name, "")
}

func (a *API) getVMParameters(name, userdata, sshkey string, ip *network.PublicIPAddress, nic *network.Interface) compute.VirtualMachine {
	osProfile := compute.OSProfile{
		AdminUsername: util.StrToPtr("core"),
		ComputerName: &name,
		LinuxConfiguration: &compute.LinuxConfiguration{
			SSH: &compute.SSHConfiguration{
				PublicKeys: &[]compute.SSHPublicKey{
					{
						Path: util.StrToPtr("/home/core/.ssh/authorized_keys"),
						KeyData: &sshkey,
					},
				},
			},
		},
	}
	if userdata != "" {
		ud := base64.StdEncoding.EncodeToString([]byte(userdata))
		osProfile.CustomData = &ud
	}
	return compute.VirtualMachine{
		Name: &name,
		Location: &a.opts.Location,
		Tags: &map[string]*string{
			"createdBy": util.StrToPtr("mantle"),
		},
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(a.opts.Size),
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: &compute.ImageReference{
					ID: &a.opts.DiskURI,
				},
				OsDisk: &compute.OSDisk{
					CreateOption: compute.FromImage,
				},
			},
			OsProfile: &osProfile,
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: nic.ID,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: util.BoolToPtr(true),
						},
					},
				},
			},
			DiagnosticsProfile: &compute.DiagnosticsProfile{
				BootDiagnostics: &compute.BootDiagnostics{
					Enabled: util.BoolToPtr(true),
					StorageURI: util.StrToPtr("https://kola.blob.core.windows.net/"),
				},
			},
		},
	}
}

func (a *API) CreateInstance(name, userdata, sshkey, resourceGroup string) (*Machine, error) {
	subnet, err := a.getSubnet(resourceGroup)
	if err != nil {
		return nil, fmt.Errorf("preparing network resources: %v", err)
	}

	ip, err := a.createPublicIP(resourceGroup)
	if err != nil {
		return nil, fmt.Errorf("creating public ip: %v", err)
	}

	nic, err := a.createNIC(ip, &subnet, resourceGroup)
	if err != nil {
		return nil, fmt.Errorf("creating nic: %v", err)
	}

	auth, err := auth.GetClientSetup(compute.DefaultBaseURI)
	if err != nil {
		return nil, err
	}
	client := compute.NewVirtualMachinesClientWithBaseURI(auth.BaseURI, auth.SubscriptionID)
	client.Authorizer = auth

	vmParams := a.getVMParameters(name, userdata, sshkey, ip, nic)

	_, err = client.CreateOrUpdate(resourceGroup, name, vmParams, nil)
	if err != nil {
		return nil, err
	}

	err = util.WaitUntilReady(5*time.Minute, 10*time.Second, func() (bool, error) {
		vm, err := client.Get(resourceGroup, name, "")
		if err != nil {
			return false, err
		}

		if vm.VirtualMachineProperties.ProvisioningState != nil && *vm.VirtualMachineProperties.ProvisioningState != "Succeeded" {
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		a.TerminateInstance(name, resourceGroup)
		return nil, fmt.Errorf("waiting for machine to become active: %v", err)
	}

	vm, err := client.Get(resourceGroup, name, "")
	if err != nil {
		return nil, err
	}

	publicaddr, privaddr, err := a.GetIPAddresses(*nic.Name, *ip.Name, resourceGroup)
	if err != nil {
		return nil, err
	}

	return &Machine{
		ID: *vm.Name,
		PublicIPAddress: publicaddr,
		PrivateIPAddress: privaddr,
		InterfaceName: *nic.Name,
		PublicIPName: *ip.Name,
	}, nil
}

func (a *API) TerminateInstance(name, resourceGroup string) error {
	auth, err := auth.GetClientSetup(compute.DefaultBaseURI)
	if err != nil {
		return err
	}
	client := compute.NewVirtualMachinesClientWithBaseURI(auth.BaseURI, auth.SubscriptionID)
	client.Authorizer = auth
	_, err = client.Delete(resourceGroup, name, nil)
	return err
}

func (a *API) GetConsoleOutput(name, resourceGroup string) ([]byte, error) {
	auth, err := auth.GetClientSetup(compute.DefaultBaseURI)
	if err != nil {
		return nil, err
	}
	client := compute.NewVirtualMachinesClientWithBaseURI(auth.BaseURI, auth.SubscriptionID)
	client.Authorizer = auth

	vm, err := client.Get(resourceGroup, name, compute.InstanceView)
	if err != nil {
		return nil, err
	}

	consoleURI := vm.VirtualMachineProperties.InstanceView.BootDiagnostics.SerialConsoleLogBlobURI
	if consoleURI == nil {
		return nil, fmt.Errorf("serial console URI is nil")
	}

	resp, err := http.Get(*consoleURI)
	if err != nil {
		return nil, fmt.Errorf("fetching console data: %v", err)
	}

	return ioutil.ReadAll(resp.Body)
}
