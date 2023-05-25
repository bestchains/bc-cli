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

package depository

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bestchains/bc-cli/pkg/account"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/nonce"
	uhttp "github.com/bestchains/bc-cli/pkg/utils/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCreateDepositoryCmd() *cobra.Command {
	var err error
	cmd := &cobra.Command{
		Use: "depository [args]",

		RunE: func(cmd *cobra.Command, args []string) error {
			// Get depot info from flags
			n, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			t, err := cmd.Flags().GetString("contentType")
			if err != nil {
				return err
			}
			id, err := cmd.Flags().GetString("contentID")
			if err != nil {
				return err
			}
			p, err := cmd.Flags().GetString("platform")
			if err != nil {
				return err
			}
			walletDir, err := cmd.Flags().GetString("wallet")
			if err != nil {
				return err
			}
			accountAddress, err := cmd.Flags().GetString("account")
			if err != nil {
				return err
			}

			// Bind the depository server host flag to viper config
			_ = viper.BindPFlag("saas.depository.server", cmd.Flags().Lookup("host"))
			host := viper.GetString("saas.depository.server")
			if host == "" {
				return fmt.Errorf("no host provided")
			}

			valueBase64 := generateValueDepotBase64(n, t, id, p)

			if accountAddress == "" {
				fmt.Println("creating untrusted depository without account endorsement")
				postValue := url.Values{}
				postValue.Add("value", valueBase64)

				// -> http://localhost/basic/putUntrustValue
				host = fmt.Sprintf("%s%s", host, common.CreateUntrustedDepository)
				resp, err := uhttp.Do(host, http.MethodPost, map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				}, []byte(postValue.Encode()))
				if err != nil {
					return err
				}
				fmt.Print(string(resp))
				return nil
			} else {
				fmt.Printf("creating trusted depository with account %s endorsement \n", accountAddress)
				//read account info
				wallet, err := account.NewLocalWallet(walletDir)
				if err != nil {
					return err
				}
				acc, err := wallet.GetAccount(accountAddress)
				if err != nil {
					return err
				}
				// get nonce
				currNonce, err := nonce.Get(host, common.DepositoryCurrentNonce, acc.Address)
				if err != nil {
					return err
				}
				// generate message
				msgBase64, err := acc.GenerateAndSignMessage(currNonce, valueBase64)
				if err != nil {
					return err
				}
				// POST PutValue request
				postValue := url.Values{}
				postValue.Add("message", msgBase64)
				postValue.Add("value", valueBase64)

				// -> http://localhost/basic/putValue
				host = fmt.Sprintf("%s%s", host, common.CreateDepository)
				resp, err := uhttp.Do(host, http.MethodPost, map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				}, []byte(postValue.Encode()))
				if err != nil {
					return err
				}
				fmt.Print(string(resp))
				return nil
			}
		},
	}
	// Set up command line flags for depository
	cmd.Flags().StringP("host", "", "", "host URL of depository server")
	cmd.Flags().StringP("wallet", "w", common.DefaultWalletConfigDir, "wallet path")
	cmd.Flags().StringP("account", "a", "", "account to be used")
	// Depository related info
	cmd.Flags().String("name", "", "depository name")
	cmd.Flags().String("contentType", "File", "depository file type")
	cmd.Flags().String("contentID", "", "depository file ID")
	cmd.Flags().String("platform", "bestchains", "depository source platform")

	// Mark the name and contentID flags as required
	err = cmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.MarkFlagRequired("contentID")
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

// generateValueDepotBase64 generates a Base64-encoded string representation of a ValueDepository object.
// The ValueDepository object contains the name, content type, content ID, and trusted timestamp of a value.
// It also includes the platform on which the value was created.
//
// Parameters:
// name (string): The name of the value.
// contentType (string): The content type of the value.
// contentID (string): The content ID of the value.
// platform (string): The platform on which the value was created.
//
// Returns:
// string: A Base64-encoded string representation of the ValueDepository object.
func generateValueDepotBase64(name string, contentType string, contentID string, platform string) string {
	// Generate ValueDepository
	valDep := ValueDepository{
		Name:             name,
		ContentType:      contentType,
		ContentID:        contentID,
		TrustedTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Platform:         platform,
	}

	// Marshal & encoding
	rawVal, err := json.Marshal(valDep)
	if err != nil {
		return ""
	}

	value := base64.StdEncoding.EncodeToString(rawVal)

	return value
}
