// Copyright 2019 Red Hat Inc.
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

package ibmcloud

import (
	"fmt"
	"os"

	"github.com/coreos/mantle/cli"
	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/api/ibmcloud"
	"github.com/coreos/pkg/capnslog"
	"github.com/spf13/cobra"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "ore/ibmcloud")

	IBMCloud = &cobra.Command{
		Use:   "ibmcloud [command]",
		Short: "ibmcloud image and vm utilities",
	}

	API             *ibmcloud.API
	region          string
	credentialsFile string
	profileName     string
	accessKeyID     string
	secretAccessKey string
)

func init() {
	defaultRegion := os.Getenv("AWS_REGION")
	if defaultRegion == "" {
		defaultRegion = "us-west-2"
	}

	// IBMCloud COS uses AWS credentials...
	IBMCloud.PersistentFlags().StringVar(&credentialsFile, "credentials-file", "", "AWS credentials file")
	IBMCloud.PersistentFlags().StringVar(&profileName, "profile", "", "AWS profile name")
	IBMCloud.PersistentFlags().StringVar(&accessKeyID, "access-id", "", "AWS access key")
	IBMCloud.PersistentFlags().StringVar(&secretAccessKey, "secret-key", "", "AWS secret key")
	IBMCloud.PersistentFlags().StringVar(&region, "region", defaultRegion, "AWS region")
	cli.WrapPreRun(IBMCloud, preflightCheck)
}

func preflightCheck(cmd *cobra.Command, args []string) error {
	api, err := ibmcloud.New(&ibmcloud.Options{
		Region:          region,
		CredentialsFile: credentialsFile,
		Profile:         profileName,
		Options:         &platform.Options{},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create IBMCloud client: %v\n", err)
		os.Exit(1)
	}
	
	API = api
	return nil
}
