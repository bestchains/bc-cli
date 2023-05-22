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

package connProfile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/bestchains/bc-cli/pkg/common"
)

func getNestedString(obj map[string]interface{}, fields ...string) string {
	val, _, _ := unstructured.NestedString(obj, fields...)
	return val
}

func NewGetConnProfileCmd(option common.Options) *cobra.Command {
	var (
		channel        string
		organization   string
		peer           string
		network        string
		id             string
		output         string
		connProfileDir string
	)

	cmd := &cobra.Command{
		Use:   "connProfile",
		Short: "Get channel's connection profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			connProfileDir = strings.TrimSuffix(connProfileDir, "/")
			_, err := os.Stat(connProfileDir)
			if err != nil {
				if !os.IsNotExist(err) {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
				err = os.MkdirAll(connProfileDir, 0755)
				if err != nil {
					fmt.Fprintln(option.ErrOut, err)
					return
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", common.InKubeGetter)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			cli, err := dynamic.NewForConfig(cfg)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			channelDetail, err := cli.Resource(schema.GroupVersionResource{Group: common.IBPGroup, Version: common.IBPVersion, Resource: common.Channel}).Get(context.TODO(), channel, v1.GetOptions{})
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			network = getNestedString(channelDetail.Object, "spec", "network")
			id = getNestedString(channelDetail.Object, "spec", "id")
			configmapName := fmt.Sprintf("chan-%s-connection-profile", channel)
			configmapDetail, err := cli.Resource(schema.GroupVersionResource{Version: common.CoreVersion, Resource: common.Configmap}).Namespace(organization).Get(context.TODO(), configmapName, v1.GetOptions{})
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			profileJson := getNestedString(configmapDetail.Object, "binaryData", "profile.json")
			rawProfileJson, err := base64.StdEncoding.DecodeString(profileJson)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}

			var profile *common.Profile
			err = json.Unmarshal(rawProfileJson, &profile)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			username := viper.GetString("auth.username")
			user := profile.Organizations[organization].Users[username]
			endpoint := profile.Peers[fmt.Sprintf("%s-%s", organization, peer)]
			var obj = make(map[string]interface{})
			obj["id"] = network
			obj["platform"] = "bestchains"
			obj["fabProfile"] = common.FabProfile{
				Channel:      id,
				Organization: organization,
				User:         user,
				Enpoint:      endpoint,
			}
			var objBytes []byte
			var targetFile string
			if output == "json" {
				objBytes, _ = json.MarshalIndent(obj, "", "  ")
				targetFile = fmt.Sprintf("%s/%s", connProfileDir, "profile.json")
			}
			if output == "yaml" {
				objBytes, _ = yaml.Marshal(obj)
				targetFile = fmt.Sprintf("%s/%s", connProfileDir, "profile.yaml")
			}
			f, err := os.Create(targetFile)
			if err != nil && !os.IsExist(err) {
				fmt.Fprintln(option.ErrOut, err)
				return
			}
			_, err = f.Write(objBytes)
			if err != nil {
				fmt.Fprintln(option.ErrOut, err)
				f.Close()
				os.Remove(targetFile)
				return
			}
			f.Close()
			fmt.Fprintf(option.Out, "connProfile %s saved\n", targetFile)
		},
	}

	cmd.Flags().StringVar(&channel, "channel", "", "channel name")
	cmd.Flags().StringVar(&organization, "org", "", "organization name")
	cmd.Flags().StringVar(&peer, "peer", "", "fabric peer name")
	cmd.Flags().StringVar(&output, "output", "json", "output file type")
	cmd.Flags().StringVar(&connProfileDir, "dir", common.DefaultConnProfileDir, "output file path")
	return cmd
}
