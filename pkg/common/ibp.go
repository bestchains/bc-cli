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

import (
	"encoding/json"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	IBPGroup             = "ibp.com"
	IBPVersion           = "v1beta1"
	IAMGroup             = "iam.tenxcloud.com"
	IAMVersion           = "v1alpha1"
	CoreVersion          = "v1"
	OrganizationResource = "organizations"
	FederationResource   = "federations"
	UserResource         = "users"
	NetworkResource      = "networks"
	Channel              = "channels"
	Configmap            = "configmaps"
	Proposal             = "proposals"
	Vote                 = "votes"
	Network              = "networks"
	EndorsePolicy        = "endorsepolicies"
	ChaincodeBuild       = "chaincodebuilds"
	Chaincode            = "chaincodes"
)

func InKubeGetter() (*clientcmdapi.Config, error) {
	kubeOIDCProxy := viper.GetString("cluster.server")
	clientID := viper.GetString("auth.clientid")
	clientSecret := viper.GetString("auth.clientsecret")
	idToken := viper.GetString("auth.idtoken")
	issuerUrl := viper.GetString("auth.issuerurl")
	refreshToken := viper.GetString("auth.refreshtoken")

	return &clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: map[string]*clientcmdapi.Cluster{
			"kube-oidc-proxy": {
				Server:                kubeOIDCProxy,
				InsecureSkipTLSVerify: true,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"oidc@kube-oidc-proxy": {
				Cluster:  "kube-oidc-proxy",
				AuthInfo: "oidc",
			},
		},
		CurrentContext: "oidc@kube-oidc-proxy",
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"oidc": {
				AuthProvider: &clientcmdapi.AuthProviderConfig{
					Name: "oidc",
					Config: map[string]string{
						// "client-id":      "bff-client",
						// "client-secret":  "61324af0-1234-4f61-b110-ef57013267d6",
						// "id-token":       "eyJhbGciOiJSUzI1NiIsImtpZCI6ImEwMzQ5ZTI0OWFjYWY4ZDQ0NzZhOGVlMzUyZGQ1YzAyNTIzMGQ2YzIifQ.eyJpc3MiOiJodHRwczovL3BvcnRhbC4xNzIuMjIuOTYuMjA5Lm5pcC5pby9vaWRjIiwic3ViIjoiQ2doaWFuZHpkMkZ1WnhJR2F6aHpZM0prIiwiYXVkIjoiYmZmLWNsaWVudCIsImV4cCI6MTY4Mzg2OTEwMywiaWF0IjoxNjgzNzgyNzAzLCJhdF9oYXNoIjoidVNjNjRyVUZiV1RMQXNWOGZtbnZzdyIsImNfaGFzaCI6InREWmJ2WGhTNTd2Uy1jY2YzUWdMMUEiLCJlbWFpbCI6ImFkbWluQHRlbnhjbG91ZC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZ3JvdXBzIjpbImlhbS50ZW54Y2xvdWQuY29tIiwib2JzZXZhYmlsaXR5IiwicmVzb3VyY2UtcmVhZGVyIiwiYmVzdGNoYWlucyJdLCJuYW1lIjoiYmp3c3dhbmciLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJiandzd2FuZyIsInBob25lIjoiIiwidXNlcmlkIjoiYmp3c3dhbmcifQ.D_lCIUT82Jg5Bje6WJSFm3fYNkaEU1bV9ZZd72qcsg2vSaI48U7Aa7tzYf6XVPrnme1TBdearuEVcDBAeRGaoEanhf_Cab1XZlsxMce7Xmny1Ih0BQE2AOE3zETNO6CeMXB8h-8tTm7UjAb8UWImTPqcmtW6VjEa5poz8_ayHpju_IYue2i6KWkTwz3AESTsPf8aa_pjk5rJqv1np7ruZjm0HJtfboQOvdbGf0u52y72N65vst431uwujKeU1sfIi5_cjMN0L4_2QG53G2eBuYswfMr9Xa7EOjK4bAUnTiuAIPzyk7P2rDOZerwkWpXZ0KVZm9cE_nOjpZDLw6AVCQ",
						// "refresh-token":  "Chl1bTc1ZndpZWs1NWRwbnA2b2NlN2ZkcnpsEhllaDJwbzRmNm83bzNwanRibGxpZHd0dzdm",
						// "idp-issuer-url": "https://portal.172.22.96.209.nip.io/oidc",
						"client-id":      clientID,
						"client-secret":  clientSecret,
						"id-token":       idToken,
						"refresh-token":  refreshToken,
						"idp-issuer-url": issuerUrl,
					},
				},
			},
		},
	}, nil
}

func GetDynamicClient() (dynamic.Interface, error) {
	cfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", InKubeGetter)
	if err != nil {
		return nil, err
	}
	cli, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func ListToObj(list corev1.List) (runtime.Object, error) {
	var obj runtime.Object
	listJSON, err := json.Marshal(list)
	if err != nil {
		return obj, err
	}
	converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, listJSON)
	if err != nil {
		return obj, err
	}
	obj = converted
	return obj, nil
}
