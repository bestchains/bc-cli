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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/printer"
	uhttp "github.com/bestchains/bc-cli/pkg/utils/http"
	"github.com/spf13/cobra"
)

func ConstructQuery(cmd *cobra.Command) string {
	query := url.Values{}
	from, _ := cmd.Flags().GetInt("from")
	if from != 0 {
		query.Add("from", fmt.Sprintf("%d", from))
	}
	size, _ := cmd.Flags().GetInt("size")
	if size != 0 {
		query.Add("size", fmt.Sprintf("%d", size))
	}
	kid, _ := cmd.Flags().GetString("kid")
	if kid != "" {
		query.Add("kid", kid)
	}
	name, _ := cmd.Flags().GetString("name")
	if name != "" {
		query.Add("name", name)
	}
	contentName, _ := cmd.Flags().GetString("contentName")
	if contentName != "" {
		query.Add("contentName", contentName)
	}

	return fmt.Sprintf("%s?%s", common.ListDepository, query.Encode())
}

var headers = []string{"index", "kid", "platform", "operator", "owner", "blockNumber", "time"}

func NewGetDepositoryCmd(option common.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depository [KID]",
		Short: "Get one or more depositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			// FIXME: the host should be read from the configuration file.
			host, _ := cmd.Flags().GetString("host")
			if host == "" {
				return fmt.Errorf("no host provided")
			}
			if len(args) == 0 {
				u := fmt.Sprintf("%s%s", host, ConstructQuery(cmd))
				x, err := uhttp.Do(u, http.MethodGet, nil, nil)
				if err != nil {
					fmt.Fprintf(option.ErrOut, "Error failed to get depository: %s\n", err.Error())
					return err
				}
				var data struct {
					Data  []Depository `json:"data"`
					Count int64        `json:"count"`
				}
				if err := json.Unmarshal(x, &data); err != nil {
					fmt.Fprintf(option.ErrOut, "unmarhsal response error %s\n", err.Error())
					return err
				}
				xx := make([]printer.Printer, len(data.Data))
				for i := 0; i < len(data.Data); i++ {
					xx[i] = data.Data[i]
				}
				printer.Print(option.Out, headers, xx)
				return nil
			}
			errMsg := make([]string, 0)
			pobj := make([]printer.Printer, 0)
			for _, kid := range args {
				u := fmt.Sprintf("%s%s", host, fmt.Sprintf(common.GetDepository, kid))
				x, err := uhttp.Do(u, http.MethodGet, nil, nil)
				if err != nil {
					errMsg = append(errMsg, err.Error())
					continue
				}
				var o Depository
				if err := json.Unmarshal(x, &o); err != nil {
					errMsg = append(errMsg, err.Error())
					continue
				}
				pobj = append(pobj, o)
			}
			printer.Print(option.Out, headers, pobj)
			for _, e := range errMsg {
				fmt.Fprintln(option.ErrOut, e)
			}
			return nil
		},
	}

	cmd.Flags().IntP("from", "f", 0, "pagination")
	cmd.Flags().IntP("size", "s", 10, "pagination size")
	cmd.Flags().StringP("kid", "k", "", "search depository by kid")
	cmd.Flags().StringP("name", "n", "", "search depository by name")
	cmd.Flags().StringP("contentName", "c", "", "search depository by content name")
	cmd.Flags().StringP("host", "", "", "bc-saas server")

	return cmd
}
