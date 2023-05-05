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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bestchains-contracts/library"
	"github.com/spf13/cobra"
)

func NewCreateAccountCmd(option common.Options) *cobra.Command {
	var (
		pkFile, walletDir string
	)
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Create an account",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if walletDir == "" {
				walletDir = common.WalletConfigDir
			}
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				pk      *ecdsa.PrivateKey
				err     error
				pkBytes []byte
			)

			if pkFile == "" {
				pk, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			} else {
				pkBytes, err = os.ReadFile(pkFile)
				if err != nil {
					return err
				}
				pkEncoded, _ := pem.Decode(pkBytes)
				pk, err = x509.ParseECPrivateKey(pkEncoded.Bytes)
			}
			if err != nil {
				return err
			}

			addr := new(library.Address)
			if err = addr.FromPublicKey(&pk.PublicKey); err != nil {
				return err
			}
			x509Encoded, err := x509.MarshalECPrivateKey(pk)
			if err != nil {
				return err
			}

			pkBytes = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
			obj := common.WalletConfig{
				Address:    addr.String(),
				PrivateKey: pkBytes,
			}

			objBytes, _ := json.Marshal(obj)
			targetFile := fmt.Sprintf("%s/%s", walletDir, addr)
			f, err := os.Create(targetFile)
			if err != nil {
				if os.IsExist(err) {
					return nil
				}
				return err
			}

			_, err = f.Write(objBytes)
			if err != nil {
				fmt.Fprintf(option.ErrOut, "Error: account/%s %s", addr, err)
				f.Close()
				os.Remove(targetFile)
				return err
			}
			f.Close()
			fmt.Fprintf(option.Out, "account/%s created\n", addr)
			return nil
		},
	}

	cmd.Flags().StringVar(&pkFile, "pk", "", "the user's own private key, which is automatically generated if not provided")
	cmd.Flags().StringVar(&walletDir, "wallet", "", "wallet path")
	return cmd
}
