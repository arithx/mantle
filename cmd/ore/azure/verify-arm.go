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
	"github.com/spf13/cobra"
)

var (
	cmdVerifyARM = &cobra.Command{
		Use:   "verify-arm",
		Short: "Verify account has ARM credentials",
		Run:   runVerifyARM,
	}

	resourceGroup  string
	storageaccount string
)

func init() {
	cmdVerifyARM.Flags().StringVar(&resourceGroup, "resource-group", "", "resource group")
	cmdVerifyARM.Flags().StringVar(&storageaccount, "storage-account", "", "storage account")

	Azure.AddCommand(cmdVerifyARM)
}

func runVerifyARM(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		plog.Fatalf("Unrecognized args in azure verify-arm: %v", args)
	}

	keys, err := api.GetStorageServiceKeysARM(resourceGroup, storageaccount)
	if err != nil {
		plog.Fatalf("fetching storage service keys: %v", err)
	}

	if keys.Keys == nil {
		plog.Fatalf("Keys is nil")
	}

	for _, key := range *keys.Keys {
		plog.Printf("%v", key.KeyName)
	}
}
