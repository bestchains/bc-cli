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
	"github.com/bestchains/bc-cli/pkg/printer"
	"github.com/spf13/cobra"
)

func NewGetAccountCmd(option common.Options) *cobra.Command {
	var (
		walletDir     string
		accountHeader = []string{"ACCOUNT"}
	)
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Display account information according to wallet path",
		PreRun: func(cmd *cobra.Command, args []string) {
			if walletDir == "" {
				walletDir = common.WalletConfigDir
			}
			walletDir = strings.TrimSuffix(walletDir, "/")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Stat(walletDir)
			if err != nil {
				fmt.Fprintf(option.ErrOut, "Error: %s\n", err)
				return nil
			}

			print := make([]printer.Printer, 0)
			dirEntries, err := os.ReadDir(walletDir)
			if err != nil {
				return err
			}
			for _, info := range dirEntries {
				if info.IsDir() {
					continue
				}
				print = append(print, AccountPrinter(info.Name()))
			}

			printer.Print(option.Out, accountHeader, print)
			return nil
		},
	}

	cmd.Flags().StringVar(&walletDir, "wallet", "", "wallet path")
	return cmd
}
