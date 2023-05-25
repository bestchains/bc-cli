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

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type Config struct {
	Auth    AuthConfig    `mapstructure:"auth"`
	Saas    SaasConfig    `mapstructure:"saas"`
	Cluster ClusterConfig `mapstructure:"cluster"`
}

// copy from "k8s.io/client-go/tools/clientcmd/api"
// Cluster contains information about how to communicate with a kubernetes cluster
type ClusterConfig struct {
	// LocationOfOrigin indicates where this object came from.  It is used for round tripping config post-merge, but never serialized.
	// +k8s:conversion-gen=false
	LocationOfOrigin string
	// Server is the address of the kubernetes cluster (https://hostname:port).
	Server string `mapstructure:"server"`
	// TLSServerName is used to check server certificate. If TLSServerName is empty, the hostname used to contact the server is used.
	// +optional
	TLSServerName string `mapstructure:"tls-server-name,omitempty"`
	// InsecureSkipTLSVerify skips the validity check for the server's certificate. This will make your HTTPS connections insecure.
	// +optional
	InsecureSkipTLSVerify bool `mapstructure:"insecure-skip-tls-verify,omitempty"`
	// CertificateAuthority is the path to a cert file for the certificate authority.
	// +optional
	CertificateAuthority string `mapstructure:"certificate-authority,omitempty"`
	// CertificateAuthorityData contains PEM-encoded certificate authority certificates. Overrides CertificateAuthority
	// +optional
	CertificateAuthorityData []byte `mapstructure:"certificate-authority-data,omitempty"`
	// ProxyURL is the URL to the proxy to be used for all requests made by this
	// client. URLs with "http", "https", and "socks5" schemes are supported.  If
	// this configuration is not provided or the empty string, the client
	// attempts to construct a proxy configuration from http_proxy and
	// https_proxy environment variables. If these environment variables are not
	// set, the client does not attempt to proxy requests.
	//
	// socks5 proxying does not currently support spdy streaming endpoints (exec,
	// attach, port forward).
	// +optional
	ProxyURL string `mapstructure:"proxy-url,omitempty"`
	// DisableCompression allows client to opt-out of response compression for all requests to the server. This is useful
	// to speed up requests (specifically lists) when client-server network bandwidth is ample, by saving time on
	// compression (server-side) and decompression (client-side): https://github.com/kubernetes/kubernetes/issues/112296.
	// +optional
	DisableCompression bool `mapstructure:"disable-compression,omitempty"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions map[string]runtime.Object `mapstructure:"extensions,omitempty"`
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
	// Username is the preferred_username(user.spec.name, not user.metadata.name)
	Username string `mapstructure:"username"`

	ClientID     string `mapstructure:"clientid"`
	ClientSecret string `mapstructure:"clientsecret"`
}

// SaasConfig represents the configuration for a SaaS application.
type SaasConfig struct {
	// Depository represents the configuration for the depository server.
	Depository Depository `mapstructure:"depository"`
	// Market represents the configuration for the market server.
	Market Market `mapstructure:"market"`
}

// Depository represents the configuration for the depository server.
type Depository struct {
	// Server represents the URL of the depository server.
	Server string `mapstructure:"server"`
}

// Market represents the configuration for the market server.
type Market struct {
	// Server represents the URL of the market server.
	Server string `mapstructure:"server"`
}

const (
	// LocalBindPort is the local bind port,
	// If you want to change it, you have to change the configuration in the oidc-server configmap at the same time.
	LocalBindPort = "127.0.0.1:26666"
	// ConfigFilePath is the config file path and file name
	ConfigFilePath = "$HOME/.bestchains/config"
	// ConfigFileType is the config file type
	ConfigFileType = "yaml"
)
