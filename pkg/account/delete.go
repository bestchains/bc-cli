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
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/spf13/cobra"
)

// NewDeleteAccountCmd returns a new cobra command for deleting an account.
// option is used to pass in common.Options.
func NewDeleteAccountCmd(option common.Options) *cobra.Command {
	// walletDir is used to specify the wallet directory.
	var walletDir string

	// cmd is the cobra command to return.
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Delete the account according to the wallet information.",

		// RunE is the function that runs when the command is called.
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a new local wallet.
			wallet, err := NewLocalWallet(walletDir)
			if err != nil {
				return err
			}

			// Delete the specified accounts.
			err = wallet.DeleteAccounts(args...)
			if err != nil {
				return err
			}
			return nil
		},
	}

	// Set the wallet directory flag.
	cmd.Flags().StringVar(&walletDir, "wallet", common.DefaultWalletConfigDir, "wallet path")
	return cmd
}
