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

package nonce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	uhttp "github.com/bestchains/bc-cli/pkg/utils/http"
)

type nonce struct {
	Nonce uint64 `json:"nonce"`
}

func Get(host string, path string, account string) (uint64, error) {
	// Add the account parameter to the URL query
	getReqValue := url.Values{}
	getReqValue.Add("account", account)

	// Construct the final URL with the query parameters
	host = fmt.Sprintf("%s%s?%s", host, path, getReqValue.Encode())

	// Make the HTTP GET request
	resp, err := uhttp.Do(host, http.MethodGet, nil, nil)
	if err != nil {
		return 0, err
	}

	// Parse the response JSON into a nonce struct
	var n nonce
	err = json.Unmarshal(resp, &n)
	if err != nil {
		return 0, err
	}

	// Return the nonce value
	return n.Nonce, nil
}
