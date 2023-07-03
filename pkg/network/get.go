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

package network

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/kubectl/pkg/cmd/get"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/federation"
	"github.com/bestchains/bc-cli/pkg/utils"
)

func NewNetworkGetCmd(option common.Options) *cobra.Command {
	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use:   "network [NAME]",
		Short: "Get a list of network",
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
				networks, err := ListNetworks(cli)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				for i := 0; i < len(networks.Items); i++ {
					list.Items = append(list.Items, runtime.RawExtension{Object: &networks.Items[i]})
				}
			} else {
				for _, arg := range utils.RemoveDuplicateForStringSlice(args) {
					network, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Network}).Get(context.TODO(), arg, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: network})
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
	return cmd
}

// ListNetworks return a list of networks.
// Return error if any error occurs.
func ListNetworks(cli dynamic.Interface) (*unstructured.UnstructuredList, error) {
	federations, err := federation.ListFederations(cli)
	if err != nil {
		return nil, err
	}
	var networkNames []string
	for _, federation := range federations.Items {
		networks, found, err := unstructured.NestedStringSlice(federation.Object, "status", "networks")
		if !found || err != nil {
			continue
		}
		networkNames = append(networkNames, networks...)
	}
	list := &unstructured.UnstructuredList{}
	for _, networkName := range utils.RemoveDuplicateForStringSlice(networkNames) {
		network, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Network}).Get(context.TODO(), networkName, v1.GetOptions{})
		if err != nil {
			continue
		}
		list.Items = append(list.Items, *network)
	}
	return list, nil
}
