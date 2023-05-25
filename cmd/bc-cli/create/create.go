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

package create

import (
	"os"

	"github.com/bestchains/bc-cli/pkg/account"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/depository"
	marketrepo "github.com/bestchains/bc-cli/pkg/market/repository"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCreateCmd returns the Cobra command for creating resources.
func NewCreateCmd() *cobra.Command {
	// Create a new Cobra command.
	cmd := &cobra.Command{
		Use: "create",
	}

	// Add the subcommands for creating a depository and market repository.
	cmd.AddCommand(depository.NewCreateDepositoryCmd())
	cmd.AddCommand(marketrepo.NewCreateMarketRepoCmd())

	// Add the subcommand for creating an account.
	cmd.AddCommand(account.NewCreateAccountCmd(common.Options{
		// Set the IOStreams for the command to use.
		IOStreams: genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		},
	}))

	// Return the command.
	return cmd
}
