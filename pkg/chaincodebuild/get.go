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

package chaincodebuild

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

func NewCCBGetCmd(option common.Options) *cobra.Command {
	var (
		network string
		id      string
		version string
	)

	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use: "ccb [NAME]",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				err := cmd.MarkFlagRequired("network")
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
			var obj runtime.Object
			if len(args) == 0 {
				var labels []string
				for k, v := range map[string]string{
					"id":      id,
					"version": version,
					"network": network,
				} {
					if v != "" {
						labels = append(labels, fmt.Sprintf("bestchains.chaincodebuild.%s=%s", k, v))
					}
				}
				chaincodeBuilds, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.ChaincodeBuild}).List(context.TODO(), v1.ListOptions{
					LabelSelector: strings.Join(labels, ","),
				})
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				for i := 0; i < len(chaincodeBuilds.Items); i++ {
					list.Items = append(list.Items, runtime.RawExtension{Object: &chaincodeBuilds.Items[i]})
				}
			} else {
				for _, arg := range args {
					chaincodeBuild, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.ChaincodeBuild}).Get(context.TODO(), arg, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						return
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: chaincodeBuild})
				}
			}

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
	cmd.Flags().StringVar(&network, "network", "", "choose a blockchain network")
	cmd.Flags().StringVar(&id, "id", "", "chaincodeBuild id")
	cmd.Flags().StringVar(&version, "version", "", "chaincodeBuild version")
	return cmd
}
