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

package azure

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cmdImageCreate = &cobra.Command{
		Use:   "image-create",
		Short: "Create Azure image",
		Long:  "Create Azure image from a blob url",
		RunE:  runImageCreate,
	}

	imageName string
	blobUrl string
	resourceGroup string
)

func init() {
	sv := cmdImageCreate.Flags().StringVar

	sv(&imageName, "img-name", "", "image name")
	sv(&blobUrl, "img-blob", "", "source blob url")
	sv(&resourceGroup, "resource-group", "kola", "resource group name")

	Azure.AddCommand(cmdImageCreate)
}

func runImageCreate(cmd *cobra.Command, args []string) error {
	_, err := api.CreateImage(imageName, resourceGroup, blobUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create image: %v\n", err)
		os.Exit(1)
	}
	return nil
}
