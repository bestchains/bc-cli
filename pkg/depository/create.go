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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bestchains/bc-cli/pkg/common"
	uhttp "github.com/bestchains/bc-cli/pkg/utils/http"
	"github.com/bestchains/bestchains-contracts/library/context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bestchains/bestchains-contracts/library"
	"github.com/spf13/cobra"
)

func NewCreateDepositoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "depository NAME [args]",
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
			untrusted, err := cmd.Flags().GetBool("untrusted")
			if err != nil {
				return err
			}

			// FIXME: the host should be read from the configuration file.
			host, _ := cmd.Flags().GetString("host")
			if host == "" {
				return fmt.Errorf("no host provided")
			}

			valueBase64 := generateValueDepotBase64(n, t, id, p)

			if untrusted {
				fmt.Print("putting untrusted value...\n")
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
				// Generate Account
				privKey, addr := randAccountAndPrivateKey()
				// Get Nonce
				non := getNonce(host, addr.String())
				// generate message
				msgBase64 := generateMessageBase64(non, privKey, []byte(valueBase64))
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

	// define flags
	cmd.Flags().StringP("host", "o", "", "host URL")
	cmd.Flags().StringP("name", "n", "", "depot name")
	cmd.Flags().StringP("contentType", "t", "", "depot file type")
	cmd.Flags().StringP("contentID", "", "", "depot file ID")
	cmd.Flags().StringP("platform", "p", "", "depot source platform")
	cmd.Flags().Bool("untrusted", true, "put untrusted value")

	return cmd
}

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

func generateMessageBase64(nonce uint64, key ecdsa.PrivateKey, value []byte) string {
	msg := context.Message{
		Nonce:     nonce,
		PublicKey: "",
		Signature: "",
	}

	// Generate signature & public key
	err := msg.GenerateSignature(&key, string(value))
	if err != nil {
		return "Fatal: generate signature failed: " + err.Error()
	}

	// Marshal & Encode
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return "Fatal: marshal failed: " + err.Error()
	}

	msgStr := base64.StdEncoding.EncodeToString(msgJson)
	if err != nil {
		return "Fatal: encode failed: " + err.Error()
	}
	return msgStr
}

func randAccountAndPrivateKey() (ecdsa.PrivateKey, library.Address) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		privateKey = new(ecdsa.PrivateKey)
	}

	var addr library.Address
	err = addr.FromPublicKey(&privateKey.PublicKey)
	if err != nil {
		addr = library.ZeroAddress
	}

	return *privateKey, addr
}

func getNonce(h string, account string) uint64 {
	getReqValue := url.Values{}
	getReqValue.Add("account", account)
	host := fmt.Sprintf("%s%s?%s", h, common.CurrentNonce, getReqValue.Encode())
	resp, err := uhttp.Do(host, http.MethodGet, nil, nil)
	if err != nil {
		return 0
	}

	var n nonce
	err = json.Unmarshal(resp, &n)
	if err != nil {
		return 0
	}

	return n.Nonce
}
