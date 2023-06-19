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
	"fmt"
	"os"
	"strings"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/spf13/cobra"
)

// NewCreateAccountCmd creates a new Cobra command for creating a new account with the given options.
// It returns the created command.
func NewCreateAccountCmd(option common.Options) *cobra.Command {
	var (
		walletDir string
	)

	cmd := &cobra.Command{
		Use:   "account",
		Short: "Create an account",

		// PreRunE is a Cobra command hook that runs before the command's RunE.
		// It checks if the wallet directory exists, creates it if not, and assigns the directory to walletDir.
		PreRunE: func(cmd *cobra.Command, args []string) error {
			walletDir = strings.TrimSuffix(walletDir, "/")
			_, err := os.Stat(walletDir)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				return os.MkdirAll(walletDir, 0755)
			}
			return nil
		},

		// RunE is the Cobra command's main function.
		// It generates a new account using the provided private key or generates a new one if none is provided.
		// It then encodes the private key and writes the account object to a file in the wallet directory.
		Run: func(cmd *cobra.Command, args []string) {
			// NewLocalWallet creates a new local wallet with a directory path.
			wallet, err := NewLocalWallet(walletDir)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			// NewAccount creates a new account.
			account, err := NewAccount()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			// StoreAccount stores the account in the wallet.
			err = wallet.StoreAccount(account)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			fmt.Fprintf(option.Out, "account/%s created\n", account.Address)
		},
	}

	// Add flags to the command.
	cmd.Flags().StringVar(&walletDir, "wallet", common.DefaultWalletConfigDir, "wallet path")
	return cmd
}
