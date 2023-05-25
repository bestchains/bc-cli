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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bestchains/bc-cli/pkg/common"
)

func TestGet(t *testing.T) {
	// Set up test data.
	testAccount := "test-account"
	testNonce := uint64(1234)

	// Create a test server that returns a JSON response containing the test nonce.
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == common.DepositoryCurrentNonce {
			if err := json.NewEncoder(w).Encode(nonce{testNonce}); err != nil {
				t.Fatalf("error when encode nonce %s", err.Error())
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer testServer.Close()

	// Call the Get function with the test server as the host.
	nonceValue, err := Get(testServer.URL, common.DepositoryCurrentNonce, testAccount)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned nonce value matches the expected test nonce.
	if nonceValue != testNonce {
		t.Fatalf("expected nonce %d, but got %d", testNonce, nonceValue)
	}
}
