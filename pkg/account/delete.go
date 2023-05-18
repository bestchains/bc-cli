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

func NewDeleteAccountCmd(option common.Options) *cobra.Command {
	var walletDir string
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Delete the account according to the wallet information.",
		PreRun: func(cmd *cobra.Command, args []string) {
			walletDir = strings.TrimSuffix(walletDir, "/")
		},
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(walletDir)
			if err != nil {
				fmt.Fprintf(option.ErrOut, "Error: %s\n", err)
				return
			}
			for _, address := range args {
				targetFile := fmt.Sprintf("%s/%s", walletDir, address)
				if err := os.Remove(targetFile); err != nil {
					fmt.Fprintf(option.ErrOut, "Error: account \"%s\" %s\n", address, err)
					continue
				}
				fmt.Fprintf(option.Out, "account \"%s\" deleted\n", address)
			}
		},
	}
	cmd.Flags().StringVar(&walletDir, "wallet", common.DefaultWalletConfigDir, "wallet path")
	return cmd
}
