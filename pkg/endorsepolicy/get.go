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

package endorsepolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/cmd/get"

	"github.com/bestchains/bc-cli/pkg/common"
)

func NewGetEndorsePolicyCmd(option common.Options) *cobra.Command {
	defaultPrinter := get.NewGetPrintFlags()

	var (
		network string
		channel string
	)
	cmd := &cobra.Command{
		Use: "ep [NAME]",
		Long: `Get endorsepolicy according to blockchain network

Examples:
  # Get all the endorsepolicies under a network
  bc-cli get ep --netowrk=<network-name>

  # Get the endorsepolicy for a number of channels specified by a network
  bc-cli get ep --network=<network-name> --channel=<channel1>,<channel2>

  # Specify the endorsepolicy name
  bc-cli get ep --network=<netowkr-name> ep1 ep2
`,

		Run: func(cmd *cobra.Command, args []string) {
			client, err := common.GetDynamicClient()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			ibpNetwork, err := client.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Network}).
				Get(context.TODO(), network, v1.GetOptions{})
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			channels, _, _ := unstructured.NestedStringSlice(ibpNetwork.Object, "status", "channels")
			if len(channels) == 0 {
				fmt.Fprintf(option.Out, "network %s don't have any channel", network)
				return
			}

			chanMap := make(map[string]struct{})
			for _, ch := range channels {
				chanMap[ch] = struct{}{}
			}

			errOutput := make([]error, 0)
			if len(channel) > 0 {
				chans := strings.Split(channel, ",")
				filterChan := make(map[string]struct{})
				for _, ch := range chans {
					if _, ok := chanMap[ch]; !ok {
						errOutput = append(errOutput, fmt.Errorf("channel %s don't belong to netowrk %s", ch, network))
						continue
					}
					filterChan[ch] = struct{}{}
				}
				if len(filterChan) == 0 {
					fmt.Fprintln(option.Out, "after filtering, there are no compliant channels.")
					return
				}
				chanMap = filterChan
			}

			epMap := make(map[string]struct{})
			if len(args) > 0 {
				for _, epName := range args {
					epMap[epName] = struct{}{}
				}
			}

			endorsePolicyList, err := client.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.EndorsePolicy}).
				List(cmd.Context(), v1.ListOptions{})
			if err != nil && !errors.IsNotFound(err) {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			if len(endorsePolicyList.Items) == 0 {
				fmt.Fprintln(option.Out, "no endorsepolicy found")
				return
			}

			outputList := corev1.List{
				TypeMeta: v1.TypeMeta{
					Kind:       "List",
					APIVersion: "v1",
				},
				ListMeta: v1.ListMeta{},
			}

			for i, item := range endorsePolicyList.Items {
				epName := item.GetName()
				epChannel, _, _ := unstructured.NestedString(item.Object, "spec", "channel")
				if _, inChannel := chanMap[epChannel]; !inChannel {
					continue
				}
				if len(epMap) > 0 {
					if _, ok := epMap[epName]; !ok {
						continue
					}
				}
				outputList.Items = append(outputList.Items, runtime.RawExtension{Object: &endorsePolicyList.Items[i]})
			}

			var obj runtime.Object

			if len(outputList.Items) != 1 {
				listData, err := json.Marshal(outputList)
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
				obj = outputList.Items[0].Object
			}

			p, err := defaultPrinter.ToPrinter()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			_ = p.PrintObj(obj, option.Out)
			for _, e := range errOutput {
				fmt.Fprintln(option.ErrOut, e)
			}
		},
	}

	defaultPrinter.AddFlags(cmd)
	cmd.Flags().StringVar(&network, "network", "", "choose a blockchain network")
	cmd.Flags().StringVar(&channel, "channel", "", "support multiple channel filtering, separated by commas")
	_ = cmd.MarkFlagRequired("network")
	return cmd
}
