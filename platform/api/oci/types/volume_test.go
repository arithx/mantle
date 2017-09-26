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

func TestCreateVolumeInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input CreateVolumeInput, jsonStr string) {
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
	displayName := "MyCustomVolume"
	sizeInMBs := 2048
	volumeBackupID := "ocid1.volume.oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"MyCustomVolume","sizeInMBs":2048,"volumeBackupId":"ocid1.volume.oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		input := CreateVolumeInput{
			AvailabilityDomain: availabilityDomain,
			CompartmentID: compartmentID,
			DisplayName: &displayName,
			SizeInMBs: &sizeInMBs,
			VolumeBackupID: &volumeBackupID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Only Required", func(t *testing.T) {
		jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"}`

		input := CreateVolumeInput{
			AvailabilityDomain: availabilityDomain,
			CompartmentID: compartmentID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestAttachVolumeInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input AttachVolumeInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	displayName := "MyCustomVolume"
	instanceID := "ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"
	volumeType := "iscsi"
	volumeID := "ocid1.volume.oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"displayName":"MyCustomVolume","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","type":"iscsi","volumeId":"ocid1.volume.oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		input := AttachVolumeInput{
			DisplayName: &displayName,
			InstanceID: instanceID,
			Type: volumeType,
			VolumeID: volumeID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Only Required", func(t *testing.T) {
		jsonStr := `{"instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","type":"iscsi","volumeId":"ocid1.volume.oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		input := AttachVolumeInput{
			InstanceID: instanceID,
			Type: volumeType,
			VolumeID: volumeID,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestVolume(t *testing.T) {
	availabilityDomain := "Uocm:PHX-AD-1"
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "MyCustomVolume"
	id := "ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"
	lifecycleState := "PROVISIONING"
	sizeInMBs := 2048
	timeCreated := "2017-09-22T21:29:30.600Z"

	requiredChecks := func(t *testing.T, output Volume) {
		strChecker(t, "Availability Domain", availabilityDomain, output.AvailabilityDomain)
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "Display Name", displayName, output.DisplayName)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		intChecker(t, "Size in MBs", sizeInMBs, output.SizeInMBs)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
	}

	jsonStr := `{"availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"MyCustomVolume","id":"ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"PROVISIONING","sizeInMBs":2048,"timeCreated":"2017-09-22T21:29:30.600Z"}`

	output := Volume{}
	err := json.Unmarshal([]byte(jsonStr), &output)
	if err != nil {
		t.Fatal(err)
	}

	requiredChecks(t, output)
}

func TestVolumeAttachment(t *testing.T) {
	attachmentType := "example-type"
	availabilityDomain := "Uocm:PHX-AD-1"
	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"
	displayName := "MyCustomVolumeAttachment"
	id := "ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"
	instanceID := "ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds"
	lifecycleState := "PROVISIONING"
	timeCreated := "2017-09-22T21:29:30.600Z"
	volumeID := "ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"

	requiredChecks:= func(t *testing.T, output VolumeAttachment) {
		strChecker(t, "Attachment Type", attachmentType, output.AttachmentType)
		strChecker(t, "Availability Domain", availabilityDomain, output.AvailabilityDomain)
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Instance ID", instanceID, output.InstanceID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
		strChecker(t, "Volume ID", volumeID, output.VolumeID)
	}

	t.Run("All Fields", func(t *testing.T) {
		jsonStr := `{"attachmentType":"example-type","availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"MyCustomVolumeAttachment","id":"ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"PROVISIONING","timeCreated":"2017-09-22T21:29:30.600Z","volumeId":"ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		output := VolumeAttachment{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
	})

	t.Run("Required Only", func(t *testing.T) {
		jsonStr := `{"attachmentType":"example-type","availabilityDomain":"Uocm:PHX-AD-1","compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","id":"ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds","instanceId":"ocid1.instance.oc1..aaaaaaaayzfqeibduyox6iib3olcmjsdlfjasldfjasldfjasdlfjaasdfds","lifecycleState":"PROVISIONING","timeCreated":"2017-09-22T21:29:30.600Z","volumeId":"ocid1.volume..oc1..fkldasjflkasjdflkasjfldsajflajsdlfjasldfjasldfjasdlfjaasdfds"}`

		output := VolumeAttachment{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Display Name", nil, output.DisplayName)
	})
}
