/*
Copyright 2023 The Bestchains Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"os"
	"path/filepath"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Constants for file paths
const (
	WalletHomeDir  = ".bestchains/wallet"      // directory for wallet files
	ConnProfileDir = ".bestchains/connProfile" // directory for connection profile files
)

// Constants for API endpoints
const (
	// Endpoint to create a depository
	CreateDepository = "/basic/putValue"
	// Endpoint to create an untrusted depository
	CreateUntrustedDepository = "/basic/putUntrustValue"
	// Endpoint to get a specific depository
	GetDepository = "/basic/depositories/%s"
	// Endpoint to list all depositories
	ListDepository = "/basic/depositories"
	// Endpoint to get the current nonce
	DepositoryCurrentNonce = "/basic/currentNonce"
	// Endpoint to download the depository certificate
	DepositoryCertificate = "/basic/depositories/certificate/%s"

	// Endpoint to create a repository
	CreateRepository = "/market/repo"
	// Endpoint to list all repositories
	ListRepositories = "/market/repos"
	// Endpoint to get the current market nonce
	MarketCurrentNonce = "/market/nonce"
)

// Variables for default directory paths
var (
	DefaultWalletConfigDir = filepath.Join(os.Getenv("HOME"), WalletHomeDir)
	DefaultConnProfileDir  = filepath.Join(os.Getenv("HOME"), ConnProfileDir)
)

// Options represents the command line options for the application
type Options struct {
	genericclioptions.IOStreams
}
