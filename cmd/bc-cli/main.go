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

package main

import (
	"bytes"
	"encoding/json"
	goflags "flag"
	"path"

	"github.com/bestchains/bc-cli/cmd/bc-cli/create"
	delcmd "github.com/bestchains/bc-cli/cmd/bc-cli/delete"
	"github.com/bestchains/bc-cli/cmd/bc-cli/get"
	"github.com/bestchains/bc-cli/pkg/auth"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bc-cli [usage]",
		Short: "Command line tools for Bestchains",
		Long:  `bc-cli is a command tool for bestchain that can query and create a variety of blockchain resources.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	fs := goflags.NewFlagSet("", goflags.PanicOnError)
	klog.InitFlags(fs)
	cmd.PersistentFlags().AddGoFlagSet(fs)
	cmd.PersistentFlags().String("issuer-url", "https://portal.172.22.96.209.nip.io/oidc", "issuer url for oidc")
	cmd.PersistentFlags().Bool("enable-auth", false, "enable oidc auth")
	cmd.PersistentFlags().String("master-url", "https://172.22.96.146:9443", "master url")
	cmd.PersistentFlags().String("client-id", "bff-client", "oidc client id")
	cmd.PersistentFlags().String("client-secret", "61324af0-1234-4f61-b110-ef57013267d6", "oidc client secret")

	ConfigFileFullPath := cmd.PersistentFlags().String("config", common.ConfigFilePath, "config file")
	_ = viper.BindPFlag("auth.issuerurl", cmd.PersistentFlags().Lookup("issuer-url"))
	_ = viper.BindPFlag("auth.enable", cmd.PersistentFlags().Lookup("enable-auth"))
	_ = viper.BindPFlag("cluster.server", cmd.PersistentFlags().Lookup("master-url"))
	_ = viper.BindPFlag("auth.clientid", cmd.PersistentFlags().Lookup("client-id"))
	_ = viper.BindPFlag("auth.clientsecret", cmd.PersistentFlags().Lookup("client-secret"))

	var config *common.Config
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		config, err = loadConfig(*ConfigFileFullPath)
		if err != nil {
			return err
		}
		configGet, err := auth.Auth(cmd.Context(), &config.Auth)
		if err != nil {
			return err
		}
		config.Auth = *configGet
		viper.Set("auth.idtoken", config.Auth.IDToken)
		viper.Set("auth.refreshtoken", config.Auth.RefreshToken)
		return nil
	}

	cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) (err error) {
		configByte, err := json.Marshal(config)
		if err != nil {
			return err
		}
		if err := viper.ReadConfig(bytes.NewReader(configByte)); err != nil {
			return err
		}
		if err := viper.WriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return viper.SafeWriteConfig()
			}
			return err
		}
		return nil
	}

	cmd.AddCommand(create.NewCreateCmd())
	cmd.AddCommand(get.NewGetCmd())
	cmd.AddCommand(delcmd.NewDeleteCmd())
	return cmd
}

func main() {
	if err := NewCmd().Execute(); err != nil {
		panic(err)
	}
}

func loadConfig(configFile string) (config *common.Config, err error) {
	config = &common.Config{}
	viper.AddConfigPath(path.Dir(configFile))
	viper.SetConfigName(path.Base(configFile))
	viper.SetConfigType(common.ConfigFileType)
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// ignore config file not exist error
			return nil, err
		}
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	klog.V(3).Infof("all config: %+v", config)
	return config, nil
}
