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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/kubectl/pkg/cmd/get"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/utils"
)

func NewFedGetCmd(option common.Options) *cobra.Command {

	defaultPrintFlag := get.NewGetPrintFlags()
	cmd := &cobra.Command{
		Use:   "fed [FED-NAME]... [-o json/yaml] [--with-org ORG-NAME]",
		Short: "Get a list of federation",
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
					federations, err := ListFederations(cli)
					if err != nil {
						fmt.Fprintln(option.ErrOut, err)
						return err
					}
					for i := 0; i < len(federations.Items); i++ {
						list.Items = append(list.Items, runtime.RawExtension{Object: &federations.Items[i]})
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

// ListFederations return a list of federations.
// Return error if any error occurs
func ListFederations(cli dynamic.Interface) (*unstructured.UnstructuredList, error) {
	username := viper.GetString("auth.username")
	users, err := cli.Resource(schema.GroupVersionResource{Group: common.IAMGroup, Version: common.IAMVersion, Resource: common.UserResource}).List(context.TODO(), v1.ListOptions{
		LabelSelector: fmt.Sprintf("t7d.io.username=%s", username),
	})
	if err != nil {
		return nil, err
	}
	if len(users.Items) == 0 {
		return nil, errors.New("No user found.")
	}
	user := users.Items[0]
	orgList := user.GetAnnotations()["bestchains"]
	var orgs map[string]interface{}
	if err := json.Unmarshal([]byte(orgList), &orgs); err != nil {
		return nil, err
	}
	orgNames, ok := orgs["list"].(map[string]interface{})
	if !ok {
		return nil, errors.New("No organization found.")
	}
	var federationNames []string
	for org := range orgNames {
		orgObj, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.OrganizationResource}).Get(context.TODO(), org, v1.GetOptions{})
		if err != nil {
			continue
		}
		feds, found, err := unstructured.NestedStringSlice(orgObj.Object, "status", "federations")
		if !found || err != nil {
			continue
		}
		federationNames = append(federationNames, feds...)
	}
	list := &unstructured.UnstructuredList{}
	for _, federationName := range utils.RemoveDuplicateForStringSlice(federationNames) {
		federation, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.FederationResource}).Get(context.TODO(), federationName, v1.GetOptions{})
		if err != nil {
			continue
		}
		list.Items = append(list.Items, *federation)
	}
	return list, nil
}
