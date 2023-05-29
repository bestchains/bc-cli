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

package federation

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
)

func NewFedGetCmd(option common.Options) *cobra.Command {

	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use: "fed [FED-NAME]... [-o json/yaml] [--with-org ORG-NAME]",
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

			orgName, _ := cmd.Flags().GetString("with-org")

			if orgName == "" {
				if len(args) == 0 {
					// get all feds of all orgs
					// - get orgList from username
					username := viper.GetString("auth.username")
					u, err := cli.Resource(schema.GroupVersionResource{
						Group:    common.IAMGroup,
						Version:  common.IAMVersion,
						Resource: common.UserResource,
					}).Get(context.TODO(), username, v1.GetOptions{})
					if err != nil {
						return err
					}
					orgList := u.GetAnnotations()["bestchains"]
					var orgs map[string]interface{}
					if err := json.Unmarshal([]byte(orgList), &orgs); err != nil {
						return err
					}
					orgNames, ok := orgs["list"].(map[string]interface{})
					if ok {
						// get org & fetch feds
						for org := range orgNames {
							orgObj, err := cli.Resource(schema.GroupVersionResource{
								Group:    common.IBPGroup,
								Version:  common.IBPVersion,
								Resource: common.OrganizationResource,
							}).Get(context.TODO(), org, v1.GetOptions{})
							if err != nil {
								fmt.Fprintln(option.ErrOut, err)
								return err
							}

							feds, _, _ := unstructured.NestedStringSlice(orgObj.Object, "status", "federations")

							// append to result list
							for _, fedName := range feds {
								fed, err := cli.Resource(schema.GroupVersionResource{
									Group:    common.IBPGroup,
									Version:  common.IBPVersion,
									Resource: common.FederationResource,
								}).Get(context.TODO(), fedName, v1.GetOptions{})
								if err != nil {
									fmt.Fprintln(option.ErrOut, err)
									continue
								}
								list.Items = append(list.Items, runtime.RawExtension{Object: fed})
							}
						}
					} else {
						fmt.Fprintln(option.Out, "No organization found.")
						return nil
					}
				} else {
					// get specified feds by fed names
					for _, arg := range args {
						fed, err := cli.Resource(schema.GroupVersionResource{
							Group:    common.IBPGroup,
							Version:  common.IBPVersion,
							Resource: common.FederationResource,
						}).Get(context.TODO(), arg, v1.GetOptions{})
						if err != nil {
							fmt.Fprintln(option.ErrOut, err)
							continue
						}
						list.Items = append(list.Items, runtime.RawExtension{Object: fed})
					}
				}
			} else {
				// get feds of certain org
				org, err := cli.Resource(schema.GroupVersionResource{
					Group:    common.IBPGroup,
					Version:  common.IBPVersion,
					Resource: common.OrganizationResource,
				}).Get(context.TODO(), orgName, v1.GetOptions{})
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return err
				}

				feds, _, _ := unstructured.NestedStringSlice(org.Object, "status", "federations")

				// append to result list
				for _, fedName := range feds {
					fed, err := cli.Resource(schema.GroupVersionResource{
						Group:    common.IBPGroup,
						Version:  common.IBPVersion,
						Resource: common.FederationResource,
					}).Get(context.TODO(), fedName, v1.GetOptions{})
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						continue
					}
					list.Items = append(list.Items, runtime.RawExtension{Object: fed})
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

	cmd.Flags().String("with-org", "", "specified organization to query for federation")
	defaultPrintFlag.AddFlags(cmd)

	return cmd
}
