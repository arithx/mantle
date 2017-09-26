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

package types

import (
        "encoding/json"
        "testing"
)

func TestVNICDetails(t *testing.T) {
	validateMarshal := func(t *testing.T, input VNICDetails, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	assignPublicIP := false
	displayName := "example-vnic-details"
	hostnameLabel := "bminstance-1"
	privateIP := "10.0.3.3"
	skipSourceDestCheck := true
	subnetID := "example-subnet-id"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"assignPublicIp":false,"displayName":"example-vnic-details","hostnameLabel":"bminstance-1","privateIp":"10.0.3.3","skipSourceDestCheck":true,"subnetId":"example-subnet-id"}`

		input := VNICDetails{
			AssignPublicIP: &assignPublicIP,
			DisplayName: &displayName,
			HostnameLabel: &hostnameLabel,
			PrivateIP: &privateIP,
			SkipSourceDestCheck: &skipSourceDestCheck,
			SubnetID: subnetID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Only Required", func(t *testing.T) {
		jsonStr := `{"subnetId":"example-subnet-id"}`

		input := VNICDetails{
			SubnetID: subnetID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestCreatePrivateIPInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input CreatePrivateIPInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	displayName := "example-private-ip"
	hostnameLabel := "bminstance-1"
	ipAddress := "10.0.3.3"
	vnicID := "ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"displayName":"example-private-ip","hostnameLabel":"bminstance-1","ipAddress":"10.0.3.3","vnicId":"ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"}`

		input := CreatePrivateIPInput{
			DisplayName: &displayName,
			HostnameLabel: &hostnameLabel,
			IPAddress: &ipAddress,
			VNICID: vnicID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{"vnicId":"ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"}`

		input := CreatePrivateIPInput{
			VNICID: vnicID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestUpdatePrivateIPInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input UpdatePrivateIPInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	displayName := "example-private-ip"
	hostnameLabel := "bminstance-1"
	vnicID := "ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"displayName":"example-private-ip","hostnameLabel":"bminstance-1","vnicId":"ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"}`

		input := UpdatePrivateIPInput{
			DisplayName: &displayName,
			HostnameLabel: &hostnameLabel,
			VNICID: &vnicID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{}`

		input := UpdatePrivateIPInput{}

		validateMarshal(t, input, jsonStr)
	})
}

func TestCreateInternetGatewayInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input CreateInternetGatewayInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-internet-gateway"
	isEnabled := true
	vnicID := "ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-internet-gateway","isEnabled":true,"vnicId":"ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"}`

		input := CreateInternetGatewayInput{
			CompartmentID: compartmentID,
			DisplayName: &displayName,
			IsEnabled: isEnabled,
			VNICID: vnicID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","isEnabled":true,"vnicId":"ocid1.example.oc1..aslkdfjaskldfjalksfjdlkasjfdsalfjdsalfjldsjafl"}`

		input := CreateInternetGatewayInput{
			CompartmentID: compartmentID,
			IsEnabled: isEnabled,
			VNICID: vnicID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestUpdateInternetGatewayInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input UpdateInternetGatewayInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	displayName := "example-internet-gateway"
	isEnabled := true

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"displayName":"example-internet-gateway","isEnabled":true}`

		input := UpdateInternetGatewayInput{
			DisplayName: &displayName,
			IsEnabled: &isEnabled,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{}`

		input := UpdateInternetGatewayInput{}

		validateMarshal(t, input, jsonStr)
	})
}

func TestAttachVNICInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input AttachVNICInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	assignPublicIP := false
	detailsDisplayName := "example-vnic-details"
	hostnameLabel := "bminstance-1"
	privateIP := "10.0.3.3"
	skipSourceDestCheck := true
	subnetID := "example-subnet-id"

	vnicDetails := VNICDetails{
		AssignPublicIP: &assignPublicIP,
		DisplayName: &detailsDisplayName,
		HostnameLabel: &hostnameLabel,
		PrivateIP: &privateIP,
		SkipSourceDestCheck: &skipSourceDestCheck,
		SubnetID: subnetID,
	}

	displayName := "example-attach-vnic"
	instanceID := "ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"createVnicDetails":{"assignPublicIp":false,"displayName":"example-vnic-details","hostnameLabel":"bminstance-1","privateIp":"10.0.3.3","skipSourceDestCheck":true,"subnetId":"example-subnet-id"},"displayName":"example-attach-vnic","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"}`

		input := AttachVNICInput{
			CreateVNICDetails: vnicDetails,
			DisplayName: &displayName,
			InstanceID: instanceID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{"createVnicDetails":{"assignPublicIp":false,"displayName":"example-vnic-details","hostnameLabel":"bminstance-1","privateIp":"10.0.3.3","skipSourceDestCheck":true,"subnetId":"example-subnet-id"},"instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"}`

		input := AttachVNICInput{
			CreateVNICDetails: vnicDetails,
			InstanceID: instanceID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestUpdateVNICInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input UpdateVNICInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	displayName := "example-update-vnic"
	hostnameLabel := "bminstance-1"
	skipSourceDestCheck := true

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"displayName":"example-update-vnic","hostnameLabel":"bminstance-1","skipSourceDestCheck":true}`

		input := UpdateVNICInput{
			DisplayName: &displayName,
			HostnameLabel: &hostnameLabel,
			SkipSourceDestCheck: &skipSourceDestCheck,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{}`

		input := UpdateVNICInput{}

		validateMarshal(t, input, jsonStr)
	})
}

func TestVNIC(t *testing.T) {
	availabilityDomain := "Uocm:PHX-AD-1"
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-vnic"
	hostnameLabel := "bminstance-1"
	id := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"
	isPrimary := true
	lifecycleState := "AVAILABLE"
	macAddress := "00:00:17:B6:4D:DD"
	privateIP := "10.0.3.3"
	publicIP := "8.8.8.8"
	skipSourceDestCheck := true
	subnetID := "example-subnet-id"
	timeCreated := "2016-08-25T21:10:29.600Z"

	requiredChecks := func(t *testing.T, output VNIC) {
		strChecker(t, "Availability Domain", availabilityDomain, output.AvailabilityDomain)
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Private IP", privateIP, output.PrivateIP)
		strChecker(t, "Subnet ID", subnetID, output.SubnetID)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
	}

	t.Run("All Fields", func (t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-vnic","hostnameLabel":"bminstance-1","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","isPrimary":true,"lifecycleState":"AVAILABLE","macAddress":"00:00:17:B6:4D:DD","privateIp":"10.0.3.3","publicIp":"8.8.8.8","skipSourceDestCheck":true,"subnetId":"example-subnet-id","timeCreated":"2016-08-25T21:10:29.600Z"}`

		output := VNIC{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
		pStrChecker(t, "Hostname Label", &hostnameLabel, output.HostnameLabel)
		pBoolChecker(t, "Is Primary", &isPrimary, output.IsPrimary)
		pStrChecker(t, "MAC Address", &macAddress, output.MACAddress)
		pStrChecker(t, "Public IP", &publicIP, output.PublicIP)
		pBoolChecker(t, "Skip Source Dest Check", &skipSourceDestCheck, output.SkipSourceDestCheck)
	})

	t.Run("Required Only", func (t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"AVAILABLE","privateIp":"10.0.3.3","subnetId":"example-subnet-id","timeCreated":"2016-08-25T21:10:29.600Z"}`

		output := VNIC{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", nil, output.DisplayName)
		pStrChecker(t, "Hostname Label", nil, output.HostnameLabel)
		pBoolChecker(t, "Is Primary", nil, output.IsPrimary)
		pStrChecker(t, "MAC Address", nil, output.MACAddress)
		pStrChecker(t, "Public IP", nil, output.PublicIP)
		pBoolChecker(t, "Skip Source Dest Check", nil, output.SkipSourceDestCheck)
	})
}

func TestVNICAttachment(t *testing.T) {
	availabilityDomain := "Uocm:PHX-AD-1"
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-vnic"
	id := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"
	instanceID := "ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"
	lifecycleState := "AVAILABLE"
	subnetID := "example-subnet-id"
	timeCreated := "2016-08-25T21:10:29.600Z"
	vlanTag := 0
	vnicID := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"

	requiredChecks := func(t *testing.T, output VNICAttachment) {
		strChecker(t, "Availability Domain", availabilityDomain, output.AvailabilityDomain)
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Instance ID", instanceID, output.InstanceID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Subnet ID", subnetID, output.SubnetID)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
	}

	t.Run("All Fields", func (t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-vnic","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"AVAILABLE","subnetId":"example-subnet-id","timeCreated":"2016-08-25T21:10:29.600Z","vlanTag":0,"vnicId":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		output := VNICAttachment{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
		pIntChecker(t, "VLAN Tag", &vlanTag, output.VLANTag)
		pStrChecker(t, "VNIC ID", &vnicID, output.VNICID)
	})

	t.Run("Required Only", func (t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"AVAILABLE","subnetId":"example-subnet-id","timeCreated":"2016-08-25T21:10:29.600Z"}`

		output := VNICAttachment{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", nil, output.DisplayName)
		pIntChecker(t, "VLAN Tag", nil, output.VLANTag)
		pStrChecker(t, "VNIC ID", nil, output.VNICID)
	})
}

func TestInternetGateway(t *testing.T) {
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-vnic"
	id := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"
	isEnabled := true
	lifecycleState := "AVAILABLE"
	timeCreated := "2016-08-25T21:10:29.600Z"
	vnicID := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"

	requiredChecks := func(t *testing.T, output InternetGateway) {
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
		strChecker(t, "VNIC ID", vnicID, output.VNICID)
	}

	t.Run("All Fields", func (t *testing.T) {
		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-vnic","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","isEnabled":true,"lifecycleState":"AVAILABLE","timeCreated":"2016-08-25T21:10:29.600Z","vnicId":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		output := InternetGateway{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
		pBoolChecker(t, "Is Enabled", &isEnabled, output.IsEnabled)
	})

	t.Run("Required Only", func (t *testing.T) {
		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"AVAILABLE","timeCreated":"2016-08-25T21:10:29.600Z","vnicId":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		output := InternetGateway{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", nil, output.DisplayName)
		pBoolChecker(t, "Is Enabled", nil, output.IsEnabled)
	})
}

func TestPrivateIP(t *testing.T) {
	availabilityDomain := "Uocm:PHX-AD-1"
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-private-ip"
	hostnameLabel := "bminstance-1"
	id := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"
	ipAddress := "10.0.3.3"
	isPrimary := true
	subnetID := "example-subnet-id"
	timeCreated := "2016-08-25T21:10:29.600Z"
	vnicID := "ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"

	t.Run("All Fields", func (t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-private-ip","hostnameLabel":"bminstance-1","id":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","ipAddress":"10.0.3.3","isPrimary":true,"subnetId":"example-subnet-id","timeCreated":"2016-08-25T21:10:29.600Z","vnicId":"ocid1.vnic..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		output := PrivateIP{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		pStrChecker(t, "Availability Domain", &availabilityDomain, output.AvailabilityDomain)
		pStrChecker(t, "Compartment ID", &compartmentID, output.CompartmentID)
		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
		pStrChecker(t, "Hostname Label", &hostnameLabel, output.HostnameLabel)
		pStrChecker(t, "ID", &id, output.ID)
		pStrChecker(t, "IP Address", &ipAddress, output.IPAddress)
		pBoolChecker(t, "Is Primary", &isPrimary, output.IsPrimary)
		pStrChecker(t, "Subnet ID", &subnetID, output.SubnetID)
		pStrChecker(t, "Time Created", &timeCreated, output.TimeCreated)
		pStrChecker(t, "VNIC ID", &vnicID, output.VNICID)
	})

	t.Run("Required Only", func (t *testing.T) {
		jsonStr := `{}`

		output := PrivateIP{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		pStrChecker(t, "Availability Domain", nil, output.AvailabilityDomain)
		pStrChecker(t, "Compartment ID", nil, output.CompartmentID)
		pStrChecker(t, "Display Name", nil, output.DisplayName)
		pStrChecker(t, "Hostname Label", nil, output.HostnameLabel)
		pStrChecker(t, "ID", nil, output.ID)
		pStrChecker(t, "IP Address", nil, output.IPAddress)
		pBoolChecker(t, "Is Primary", nil, output.IsPrimary)
		pStrChecker(t, "Subnet ID", nil, output.SubnetID)
		pStrChecker(t, "Time Created", nil, output.TimeCreated)
		pStrChecker(t, "VNIC ID", nil, output.VNICID)
	})
}
