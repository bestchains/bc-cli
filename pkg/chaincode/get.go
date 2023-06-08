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

package chaincode

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/cmd/get"

	"github.com/bestchains/bc-cli/pkg/common"
)

func NewCCGetCmd(option common.Options) *cobra.Command {
	var (
		channel string
		id      string
		version string
	)
	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use:   "chaincode [NAME]",
		Short: "Get a list of the chaincode installed on a channel",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				err := cmd.MarkFlagRequired("channel")
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := common.GetDynamicClient()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			list := corev1.List{
				TypeMeta: v1.TypeMeta{
					Kind:       "List",
					APIVersion: "v1",
				},
				ListMeta: v1.ListMeta{},
			}
			if len(args) == 0 {
				var labels []string
				for k, v := range map[string]string{
					"channel": channel,
					"id":      id,
					"version": version,
				} {
					if v != "" {
						labels = append(labels, fmt.Sprintf("bestchains.chaincode.%s=%s", k, v))
					}
				}
				chaincodes, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Chaincode}).List(context.TODO(), v1.ListOptions{
					LabelSelector: strings.Join(labels, ","),
				})
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				for i := 0; i < len(chaincodes.Items); i++ {
					list.Items = append(list.Items, runtime.RawExtension{Object: &chaincodes.Items[i]})
				}
			} else {
				for _, arg := range args {
					chaincode, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Chaincode}).Get(context.TODO(), arg, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: chaincode})
				}
			}
			var obj runtime.Object
			if len(list.Items) != 1 {
				listData, err := json.Marshal(list)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, listData)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				obj = converted
			} else {
				obj = list.Items[0].Object
			}

			p, err := defaultPrintFlag.ToPrinter()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			_ = p.PrintObj(obj, option.Out)
		},
	}

	defaultPrintFlag.AddFlags(cmd)
	cmd.Flags().StringVar(&channel, "channel", "", "channel name")
	cmd.Flags().StringVar(&id, "id", "", "chaincode id")
	cmd.Flags().StringVar(&version, "version", "", "chaincode version")
	return cmd
}
