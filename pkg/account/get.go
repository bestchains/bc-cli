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

package account

import (
	"strings"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/printer"
	"github.com/spf13/cobra"
)

// NewGetAccountCmd creates a new Cobra command for displaying account information
// according to wallet path.
func NewGetAccountCmd(option common.Options) *cobra.Command {
	// Initialize variables.
	var (
		walletDir     string
		accountHeader = []string{"ACCOUNT"}
	)

	// Create the command.
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Display account information according to wallet path",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Remove trailing slash from wallet path.
			walletDir = strings.TrimSuffix(walletDir, "/")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a new local wallet.
			wallet, err := NewLocalWallet(walletDir)
			if err != nil {
				return err
			}

			// Get a list of accounts from the wallet.
			accounts, err := wallet.ListAccounts()
			if err != nil {
				return err
			}

			// Create a list of account printers.
			print := make([]printer.Printer, 0)
			for _, account := range accounts {
				print = append(print, AccountPrinter(account))
			}

			// Print the account information.
			printer.Print(option.Out, accountHeader, print)
			return nil
		},
	}

	// Add the wallet flag to the command.
	cmd.Flags().StringVar(&walletDir, "wallet", common.DefaultWalletConfigDir, "wallet path")
	return cmd
}
