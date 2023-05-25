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

package repository

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/bestchains/bc-cli/pkg/account"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/nonce"
	uhttp "github.com/bestchains/bc-cli/pkg/utils/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCreateMarketRepoCmd() *cobra.Command {
	var err error
	cmd := &cobra.Command{
		Use: "market repo [args]",
		RunE: func(cmd *cobra.Command, args []string) error {
			// WalletDir&Account will be used to init the wallet
			walletDir, err := cmd.Flags().GetString("wallet")
			if err != nil {
				return err
			}
			accountAddress, err := cmd.Flags().GetString("account")
			if err != nil {
				return err
			}

			// read repository url
			repoURL, err := cmd.Flags().GetString("url")
			if err != nil {
				return err
			}

			// bind depository server to flag
			_ = viper.BindPFlag("saas.market.server", cmd.Flags().Lookup("host"))
			host := viper.GetString("saas.market.server")
			if host == "" {
				return fmt.Errorf("no host provided")
			}

			fmt.Printf("creating repository with account %s endorsement \n", accountAddress)
			resp, err := CreateRepo(host, walletDir, accountAddress, repoURL)
			if err != nil {
				return err
			}
			fmt.Print(string(resp))
			return nil

		},
	}

	// define flags
	cmd.Flags().StringP("host", "", "", "host URL of market server")
	cmd.Flags().StringP("wallet", "w", common.DefaultWalletConfigDir, "wallet path")
	cmd.Flags().StringP("account", "a", "", "account to be used")
	cmd.Flags().String("repo-url", "", "repository url")

	// define required flags
	err = cmd.MarkFlagRequired("account")
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired("repo-url")
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

// CreateRepo creates a new repository on the specified host using the provided account and repo URL.
// It returns the response body as a byte slice and any error encountered.
func CreateRepo(host string, walletDir string, accountAddress string, repoURL string) ([]byte, error) {
	// Read account info.
	wallet, err := account.NewLocalWallet(walletDir)
	if err != nil {
		return nil, err
	}
	acc, err := wallet.GetAccount(accountAddress)
	if err != nil {
		return nil, err
	}

	// Get nonce.
	currNonce, err := nonce.Get(host, common.MarketCurrentNonce, acc.Address)
	if err != nil {
		return nil, err
	}

	// Generate message.
	msgBase64, err := acc.GenerateAndSignMessage(currNonce, repoURL)
	if err != nil {
		return nil, err
	}

	// POST PutValue request.
	postValue := url.Values{}
	postValue.Add("message", msgBase64)
	postValue.Add("url", repoURL)

	host = fmt.Sprintf("%s%s", host, common.CreateRepository)
	resp, err := uhttp.Do(host, http.MethodPost, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, []byte(postValue.Encode()))
	if err != nil {
		return nil, err
	}

	return resp, nil
}
