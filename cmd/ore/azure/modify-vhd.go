// Copyright 2020 Red Hat
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
	"io/ioutil"
	"strings"

	"github.com/Microsoft/azure-vhd-utils/vhdcore/diskstream"
	"github.com/Microsoft/azure-vhd-utils/vhdcore/validator"
	"github.com/spf13/cobra"
)

var (
	cmdModifyVHD = &cobra.Command{
		Use:   "modify-vhd",
		Short: "Modify a VHD",
		Run:   runModifyVHD,
	}

	localOut string
)

func init() {
	bv := cmdModifyVHD.Flags().BoolVar
	sv := cmdModifyVHD.Flags().StringVar

	bv(&ubo.validate, "validate", true, "validate blob as VHD file")
	sv(&ubo.vhd, "file", "", "path to CoreOS image")
	sv(&localOut, "output", "", "path to output image")

	Azure.AddCommand(cmdModifyVHD)
}

func runModifyVHD(cmd *cobra.Command, args []string) {
	if err := api.SetupClients(); err != nil {
		plog.Fatalf("setting up clients: %v\n", err)
	}

	if ubo.validate {
		plog.Printf("Validating VHD %q", ubo.vhd)
		if !strings.HasSuffix(strings.ToLower(ubo.vhd), ".vhd") {
			plog.Fatalf("Image should end with .vhd")
		}

		if err := validator.ValidateVhd(ubo.vhd); err != nil {
			plog.Fatal(err)
		}

		if err := validator.ValidateVhdSize(ubo.vhd); err != nil {
			plog.Fatal(err)
		}
	}

	ds, err := diskstream.CreateNewDiskStream(ubo.vhd)
	if err != nil {
		plog.Fatal("creating diskstream: %v", err)
	}

	buffer := make([]byte, ds.GetSize())
	_, err = ds.Read(buffer)
	if err != nil {
		plog.Fatal("reading diskstream: %v", err)
	}

	err = ioutil.WriteFile(localOut, buffer, 0644)
	if err != nil {
		plog.Fatal("writing output file: %v", err)
	}
}
