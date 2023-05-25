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

package channel

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/cmd/get"

	"github.com/bestchains/bc-cli/pkg/common"
)

// NewChanGetCmd returns the `kubectl get` command for channel.
func NewChanGetCmd(option common.Options) *cobra.Command {

	defaultPrintFlag := get.NewGetPrintFlags()

	cmd := &cobra.Command{
		Use: "channel [NAME] -n NETWORK-NAME",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get config from getter
			cli, err := common.GetDynamicClient()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return err
			}
			list := corev1.List{
				TypeMeta: v1.TypeMeta{
					Kind:       "List",
					APIVersion: "v1",
				},
				ListMeta: v1.ListMeta{},
			}
			var obj runtime.Object

			netName, _ := cmd.Flags().GetString("network")

			if len(args) == 0 {
				// list all channel of the given network
				net, err := cli.Resource(schema.GroupVersionResource{
					Group:    common.IBPGroup,
					Version:  common.IBPVersion,
					Resource: common.NetworkResource,
				}).Get(context.TODO(), netName, v1.GetOptions{})
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return err
				}

				channels, _, _ := unstructured.NestedStringSlice(net.Object, "status", "channels")

				for _, chName := range channels {
					ch, err := cli.Resource(schema.GroupVersionResource{
						Group:    common.IBPGroup,
						Version:  common.IBPVersion,
						Resource: common.Channel,
					}).Get(context.TODO(), chName, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						return err
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: ch})
				}

			} else {
				// get channel by name
				for _, arg := range args {
					ch, err := cli.Resource(schema.GroupVersionResource{
						Group:    common.IBPGroup,
						Version:  common.IBPVersion,
						Resource: common.Channel,
					}).Get(context.TODO(), arg, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: ch})
				}
			}

			// object process
			if len(list.Items) != 1 {
				obj, err = common.ListToObj(list)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return err
				}
			} else {
				obj = list.Items[0].Object
			}

			// print result
			p, err := defaultPrintFlag.ToPrinter()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return err
			}
			_ = p.PrintObj(obj, option.Out)
			return nil
		},
	}

	cmd.Flags().StringP("network", "n", "", "network of the desired channel")
	_ = cmd.MarkFlagRequired("network")

	return cmd
}
