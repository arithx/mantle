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

func TestLaunchInstanceInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input LaunchInstanceInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	availabilityDomain := "Uocm:PHX-AD-1"
        compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-instance"
	hostnameLabel := "bminstance-1"
	imageID := "ocid1.image..."
	ipxeScript := "example-ipxe-script"
	metadata := map[string]string{
		"foo": "bar",
	}
	shape := "VM.Standard1.1"
	subnetID := "example-subnet-id"
	extendedMetadata := json.RawMessage(`{"extended":{"metadata":"example"}}`)

	assignPublicIP := false
        detailsDisplayName := "example-vnic-details"
        privateIP := "10.0.3.3"
        skipSourceDestCheck := true

        vnicDetails := VNICDetails{
                AssignPublicIP: &assignPublicIP,
                DisplayName: &detailsDisplayName,
                HostnameLabel: &hostnameLabel,
                PrivateIP: &privateIP,
                SkipSourceDestCheck: &skipSourceDestCheck,
                SubnetID: subnetID,
        }

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","createVnicDetails":{"assignPublicIp":false,"displayName":"example-vnic-details","hostnameLabel":"bminstance-1","privateIp":"10.0.3.3","skipSourceDestCheck":true,"subnetId":"example-subnet-id"},"displayName":"example-instance","hostnameLabel":"bminstance-1","imageId":"ocid1.image...","ipxeScript":"example-ipxe-script","metadata":{"foo":"bar"},"shape":"VM.Standard1.1","subnetId":"example-subnet-id","extendedMetadata":{"extended":{"metadata":"example"}}}`

		input := LaunchInstanceInput{
			AvailabilityDomain: availabilityDomain,
			CompartmentID: compartmentID,
			CreateVNICDetails: &vnicDetails,
			DisplayName: &displayName,
			HostnameLabel: &hostnameLabel,
			ImageID: imageID,
			IPXEScript: &ipxeScript,
			Metadata: &metadata,
			Shape: shape,
			SubnetID: &subnetID,
			ExtendedMetadata: &extendedMetadata,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","imageId":"ocid1.image...","shape":"VM.Standard1.1"}`

		input := LaunchInstanceInput{
			AvailabilityDomain: availabilityDomain,
			CompartmentID: compartmentID,
			ImageID: imageID,
			Shape: shape,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestListInstancesInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input ListInstancesInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	availabilityDomain := "Uocm:PHX-AD-1"
        compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
        displayName := "example-instance"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-instance"}`

		input := ListInstancesInput{
			AvailabilityDomain: &availabilityDomain,
			CompartmentID: compartmentID,
			DisplayName: &displayName,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("RequiredOnly", func(t *testing.T) {
		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"}`

		input := ListInstancesInput{
			CompartmentID: compartmentID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestInstanceConsoleConnectionInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input InstanceConsoleConnectionInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	instanceID := "ocid1.instance.oc1.phx.abyhqljrqyriphyccj75yut36ybxmlfgawtl7m77vqanhg6w4bdszaitd3da"
	publicKey := "example-public-key"

	jsonStr := `{"instanceId":"ocid1.instance.oc1.phx.abyhqljrqyriphyccj75yut36ybxmlfgawtl7m77vqanhg6w4bdszaitd3da","publicKey":"example-public-key"}`

	input := InstanceConsoleConnectionInput{
		InstanceID: instanceID,
		PublicKey: publicKey,
	}

	validateMarshal(t, input, jsonStr)
}

func TestCaptureConsoleHistoryInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input CaptureConsoleHistoryInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	instanceID := "ocid1.instance.oc1.phx.abyhqljrqyriphyccj75yut36ybxmlfgawtl7m77vqanhg6w4bdszaitd3da"

	jsonStr := `{"instanceId":"ocid1.instance.oc1.phx.abyhqljrqyriphyccj75yut36ybxmlfgawtl7m77vqanhg6w4bdszaitd3da"}`

	input := CaptureConsoleHistoryInput{
		InstanceID: instanceID,
	}

	validateMarshal(t, input, jsonStr)
}

func TestShape(t *testing.T) {
	validateMarshal := func(t *testing.T, input Shape, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	shape := "VM-Standard1.1"

	jsonStr := `{"shape":"VM-Standard1.1"}`

	input := Shape{
		Shape: shape,
	}

	validateMarshal(t, input, jsonStr)
}

func TestInstanceConsoleConnection(t *testing.T) {
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	connectionString := "example-connection-string"
	fingerprint := "aa:bb:cc:dd:00:11:22:33:44"
	id := "ocid1.console.connection..."
	instanceID := "ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"
	lifecycleState := "ACTIVE"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","connectionString":"example-connection-string","fingerprint":"aa:bb:cc:dd:00:11:22:33:44","id":"ocid1.console.connection...","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"ACTIVE"}`

		output := InstanceConsoleConnection{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		pStrChecker(t, "Compartment ID", &compartmentID, output.CompartmentID)
		pStrChecker(t, "Connection String", &connectionString, output.ConnectionString)
		pStrChecker(t, "Fingerprint", &fingerprint, output.Fingerprint)
		pStrChecker(t, "ID", &id, output.ID)
		pStrChecker(t, "Instance ID", &instanceID, output.InstanceID)
		pStrChecker(t, "Lifecycle State", &lifecycleState, output.LifecycleState)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{}`

		output := InstanceConsoleConnection{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		pStrChecker(t, "Compartment ID", nil, output.CompartmentID)
		pStrChecker(t, "Connection String", nil, output.ConnectionString)
		pStrChecker(t, "Fingerprint", nil, output.Fingerprint)
		pStrChecker(t, "ID", nil, output.ID)
		pStrChecker(t, "Instance ID", nil, output.InstanceID)
		pStrChecker(t, "Lifecycle State", nil, output.LifecycleState)
	})
}

func TestInstance(t *testing.T) {
	availabilityDomain := "Uocm:PHX-AD-1"
        compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-instance"
	id := "ocid1.instance..."
	imageID := "ocid1.image..."
	ipxeScript := "example-ipxe-script"
	lifecycleState := "RUNNING"
	region := "phx"
	shape := "VM.Standard1.1"
	extendedMetadata := json.RawMessage(`{"extended":{"metadata":"example"}}`)
	timeCreated := "2016-08-25T21:10:29.600Z"

	requiredChecks := func(t *testing.T, output Instance) {
		strChecker(t, "Availability Domain", availabilityDomain, output.AvailabilityDomain)
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Region", region, output.Region)
		strChecker(t, "Shape", shape, output.Shape)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
	}

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-instance","id":"ocid1.instance...","imageId":"ocid1.image...","ipxeScript":"example-ipxe-script","lifecycleState":"RUNNING","metadata":{"foo":"bar"},"region":"phx","shape":"VM.Standard1.1","timeCreated":"2016-08-25T21:10:29.600Z","extendedMetadata":{"extended":{"metadata":"example"}}}`


		output := Instance{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
		pStrChecker(t, "Image ID", &imageID, output.ImageID)
		pStrChecker(t, "iPXE Script", &ipxeScript, output.IPXEScript)

		for key, val := range *output.Metadata {
			if key != "foo" || val != "bar" {
				t.Fatalf("Metadata is wrong")
			}
		}

		if output.ExtendedMetadata == nil {
			t.Fatalf("ExtendedMetadata is wrong")
		}
		strChecker(t, "Extended Metadata", string(extendedMetadata), string(*output.ExtendedMetadata))
	})

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","id":"ocid1.instance...","lifecycleState":"RUNNING","region":"phx","shape":"VM.Standard1.1","timeCreated":"2016-08-25T21:10:29.600Z"}`


		output := Instance{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", nil, output.DisplayName)
		pStrChecker(t, "Image ID", nil, output.ImageID)
		pStrChecker(t, "iPXE Script", nil, output.IPXEScript)

		if output.Metadata != nil {
			t.Fatalf("Metadata is wrong")
		}

		if output.ExtendedMetadata != nil {
			t.Fatalf("ExtendedMetadata is wrong")
		}
	})
}

func TestConsoleHistory(t *testing.T) {
	availabilityDomain := "Uocm:PHX-AD-1"
        compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "example-instance"
	id := "ocid1.console..."
	instanceID := "ocid1.instance..."
	lifecycleState := "RUNNING"
	timeCreated := "2016-08-25T21:10:29.600Z"

	requiredChecks := func(t *testing.T, output ConsoleHistory) {
		strChecker(t, "Availability Domain", availabilityDomain, output.AvailabilityDomain)
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Instance ID", instanceID, output.InstanceID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
	}

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"example-instance","id":"ocid1.console...","instanceId":"ocid1.instance...","lifecycleState":"RUNNING","timeCreated":"2016-08-25T21:10:29.600Z"}`

		output := ConsoleHistory{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","id":"ocid1.console...","instanceId":"ocid1.instance...","lifecycleState":"RUNNING","timeCreated":"2016-08-25T21:10:29.600Z"}`

		output := ConsoleHistory{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", nil, output.DisplayName)
	})
}
