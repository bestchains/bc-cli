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

package proposal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/cmd/get"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/org"
	"github.com/bestchains/bc-cli/pkg/utils"
)

func NewProposalGetCmd(option common.Options) *cobra.Command {
	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use:   "proposal [NAME]",
		Short: "Get a list of proposal",
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
				username := viper.GetString("auth.username")
				organizations, err := org.ListOrganizations(cli, fmt.Sprintf("bestchains.organization.admin=%s", username), "")
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				var proposalNames []string
				for _, org := range organizations.Items {
					namespace := org.GetName()
					votes, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Vote}).Namespace(namespace).List(context.TODO(), v1.ListOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					for _, vote := range votes.Items {
						proposalName := utils.GetNestedString(vote.Object, "spec", "proposalName")
						proposalNames = append(proposalNames, proposalName)
					}
				}
				for _, proposalName := range utils.RemoveDuplicateForStringSlice(proposalNames) {
					proposal, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Proposal}).Get(context.TODO(), proposalName, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: proposal})
				}
			} else {
				for _, arg := range utils.RemoveDuplicateForStringSlice(args) {
					proposal, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Proposal}).Get(context.TODO(), arg, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: proposal})
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
