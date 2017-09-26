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

func TestCreateImageInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input CreateImageInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	compartmentID := "ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq"

	t.Run("No ImageSource", func(t *testing.T) {
		instanceID := "ocid1.instance.oc1.phx.abyhqljrqyriphyccj75yut36ybxmlfgawtl7m77vqanhg6w4bdszaitd3da"
		displayName := "MyCustomImage"

		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"MyCustomImage","instanceId":"ocid1.instance.oc1.phx.abyhqljrqyriphyccj75yut36ybxmlfgawtl7m77vqanhg6w4bdszaitd3da"}`

		input := CreateImageInput{
			CompartmentID: compartmentID,
			DisplayName: &displayName,
			InstanceID: &instanceID,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Object Tuple ImageSource", func(t *testing.T) {
		imageSource := ImageSourceViaObjectStorageTuple{
			ObjectName: "image-to-import.qcow2",
			BucketName: "MyBucket",
			NamespaceName: "MyNamespace",
			SourceType: "objectStorageTuple",
		}

		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","imageSourceDetails":{"bucketName":"MyBucket","namespaceName":"MyNamespace","objectName":"image-to-import.qcow2","sourceType":"objectStorageTuple"}}`
		input := CreateImageInput{
			CompartmentID: compartmentID,
			ImageSource: &imageSource,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Object Storage Service URL", func(t *testing.T) {
		displayName := "MyImportedImage"
		imageSource := ImageSourceViaObjectStorageURI{
			SourceURI: "https://objectstorage.us-phoenix-1.oraclecloud.com/n/MyNamespace/b/MyBucket/o/image-to-import.qcow2",
			SourceType: "objectStorageUri",
		}

		jsonStr := `{"compartmentId":"ocid1.compartment.oc1..aaaaaaaayzfqeibduyox6iib3olcmdar3ugly4fmameq4h7lcdlihrvur7xq","displayName":"MyImportedImage","imageSourceDetails":{"sourceUri":"https://objectstorage.us-phoenix-1.oraclecloud.com/n/MyNamespace/b/MyBucket/o/image-to-import.qcow2","sourceType":"objectStorageUri"}}`

		input := CreateImageInput{
			CompartmentID: compartmentID,
			DisplayName: &displayName,
			ImageSource: &imageSource,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestExportImageInput(t *testing.T) {
	validateMarshal := func(t *testing.T, input ExportImageInput, jsonStr string) {
		marshal, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		if string(marshal) != jsonStr {
			t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
		}
	}

	t.Run("Namespace, Bucket Name, and Object Name", func (t *testing.T) {
		objectName := "exported-image.qcow2"
		bucketName := "MyBucket"
		namespace := "MyNamespace"
		destinationType := "objectStorageTuple"

		jsonStr := `{"destinationType":"objectStorageTuple","bucketName":"MyBucket","namespaceName":"MyNamespace","objectName":"exported-image.qcow2"}`

		input := ExportImageInput{
			DestinationType: destinationType,
			BucketName: &bucketName,
			NamespaceName: &namespace,
			ObjectName: &objectName,
		}

		validateMarshal(t, input, jsonStr)
	})

	t.Run("Object Storage URL", func (t *testing.T) {
		destinationUri := "https://objectstorage.us-phoenix-1.oraclecloud.com/n/MyNamespace/b/MyBucket/o/exported-image.qcow2"
		destinationType := "objectStorageUri"

		jsonStr := `{"destinationType":"objectStorageUri","destinationUri":"https://objectstorage.us-phoenix-1.oraclecloud.com/n/MyNamespace/b/MyBucket/o/exported-image.qcow2"}`

		input := ExportImageInput{
			DestinationType: destinationType,
			DestinationURI: &destinationUri,
		}

		validateMarshal(t, input, jsonStr)
	})
}

func TestUpdateImageInput(t *testing.T) {
	displayName := "MyFavoriteImage"

	jsonStr := `{"displayName":"MyFavoriteImage"}`

	input := UpdateImageInput{
		DisplayName: &displayName,
	}

	marshal, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	if string(marshal) != jsonStr {
		t.Fatalf("JSON doens't match:\n\tExpected: %s\n\tReceived: %s", jsonStr, marshal)
	}
}

func TestImage(t *testing.T) {
	baseImageID := "example-base-image-id"
	compartmentID := "example-compartment-id"
	createImageAllowed := true
	displayName := "example-display-name"
	id := "example-id"
	lifecycleState := "AVAILABLE"
	operatingSystem := "coreos"
	operatingSystemVersion := "1535.2.0"
	timeCreated := "2017-09-22T21:29:30.600Z"

	requiredChecks := func(t *testing.T, output Image) {
		strChecker(t, "Compartment ID", compartmentID, output.CompartmentID)
		boolChecker(t, "Create Image Allowed", createImageAllowed, output.CreateImageAllowed)
		strChecker(t, "ID", id, output.ID)
		strChecker(t, "Lifecycle State", lifecycleState, output.LifecycleState)
		strChecker(t, "Operating System", operatingSystem, output.OperatingSystem)
		strChecker(t, "Operating System Version", operatingSystemVersion, output.OperatingSystemVersion)
		strChecker(t, "Time Created", timeCreated, output.TimeCreated)
	}

	t.Run("All Fields", func (t *testing.T) {
		jsonStr := `{"baseImageId":"example-base-image-id","compartmentId":"example-compartment-id","createImageAllowed":true,"displayName":"example-display-name","id":"example-id","lifecycleState":"AVAILABLE","operatingSystem":"coreos","operatingSystemVersion":"1535.2.0","timeCreated":"2017-09-22T21:29:30.600Z"}`

		output := Image{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Base Image ID", &baseImageID, output.BaseImageID)
		pStrChecker(t, "Display Name", &displayName, output.DisplayName)
	})

	t.Run("Required Fields Only", func (t *testing.T) {
		jsonStr := `{"compartmentId":"example-compartment-id","createImageAllowed":true,"id":"example-id","lifecycleState":"AVAILABLE","operatingSystem":"coreos","operatingSystemVersion":"1535.2.0","timeCreated":"2017-09-22T21:29:30.600Z"}`

		output := Image{}
		err := json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			t.Fatal(err)
		}

		requiredChecks(t, output)

		pStrChecker(t, "Base Image ID", nil, output.BaseImageID)
		pStrChecker(t, "Display Name", nil, output.DisplayName)
	})
}
