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

package org

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
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/bestchains/bc-cli/pkg/common"
)

func NewOrgGetCmd(option common.Options) *cobra.Command {
	var (
		labelSelector string
		fieldSelector string
	)

	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use:   "org [NAME]",
		Short: "Get a list of organization",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			if len(args) == 0 {
				orgs, err := ListOrganizations(cli, labelSelector, fieldSelector)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return err
				}
				for i := 0; i < len(orgs.Items); i++ {
					list.Items = append(list.Items, runtime.RawExtension{Object: &orgs.Items[i]})
				}
			} else {
				for _, arg := range args {
					org, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.OrganizationResource}).Get(context.TODO(), arg, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: org})
				}
			}

			if len(list.Items) != 1 {
				listData, err := json.Marshal(list)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return err
				}
				converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, listData)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return err
				}
				obj = converted
			} else {
				obj = list.Items[0].Object
			}

			p, err := defaultPrintFlag.ToPrinter()
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return err
			}
			_ = p.PrintObj(obj, option.Out)
			return nil
		},
	}

	defaultPrintFlag.AddFlags(cmd)
	cmd.Flags().StringVar(&fieldSelector, "field-selector", "", "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.")
	cmdutil.AddLabelSelectorFlagVar(cmd, &labelSelector)

	return cmd
}

// ListOrganizations returns a list of organizations filtered by labelSelector and fieldSelector.
// Return error if any error occurs
func ListOrganizations(cli dynamic.Interface, labelSelector string, fieldSelector string) (*unstructured.UnstructuredList, error) {
	organizations, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.OrganizationResource}).List(context.TODO(), v1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return organizations, nil
}
