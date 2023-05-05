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

package common

// note: mapstructure tag is used for viper auto unmashal,
// viper's own restrictions can only use lowercase letters,
// in 2019 the community found this point is not very reasonable, but has not changed
// see https://github.com/spf13/viper/pull/758

type Config struct {
	Auth AuthConfig `mapstructure:"auth"`
	Saas SaasConfig `mapstructure:"saas"`
}

type AuthConfig struct {
	// Enable is the enable flag
	Enable bool `mapstructure:"enable"`
	// IssuerURL is the URL of the OIDC issuer.
	IssuerURL string `mapstructure:"issuerurl"`
	// IDToken is the id-token
	IDToken string `mapstructure:"idtoken"`
	// RefreshToken is the refresh-token
	RefreshToken string `mapstructure:"refreshtoken"`
	// Expiry is the expiry time of the access token
	Expiry int64 `mapstructure:"expiry"`
}

type SaasConfig struct {
	Depository Depository `mapstructure:"depository"`
}

type Depository struct {
	Server string `mapstructure:"server"`
}

const (
	// LocalBindPort is the local bind port,
	// If you want to change it, you have to change the configuration in the oidc-server configmap at the same time.
	LocalBindPort = "127.0.0.1:26666"
	// ClientID for oidc
	// If you want to change it, you have to change the configuration in the oidc-server configmap at the same time.
	ClientID = "bc-cli"
	// ClientSecret for oidc
	// If you want to change it, you have to change the configuration in the oidc-server configmap at the same time.
	ClientSecret = "bc-cli-cli"

	// ConfigFilePath is the config file path and file name
	ConfigFilePath = "$HOME/.bestchains/config"
	// ConfigFileType is the config file type
	ConfigFileType = "yaml"
)
