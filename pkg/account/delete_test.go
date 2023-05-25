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

package account

import (
	"bytes"
	"os"
	"testing"

	"github.com/bestchains/bc-cli/pkg/common"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestNewDeleteAccountCmd(t *testing.T) {
	// Set up test data.
	testAccounts := []string{"test-account-1", "test-account-2"}

	// Create a temporary directory for the wallet.
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Create a new LocalWallet instance.
	wallet, err := NewLocalWallet(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Store some test accounts in the wallet.
	for _, account := range testAccounts {
		err = wallet.StoreAccount(&Account{Address: account})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create a new command and execute it.
	buf := new(bytes.Buffer)
	testOpts := common.Options{
		IOStreams: genericclioptions.IOStreams{
			Out: buf,
		},
	}
	cmd := NewDeleteAccountCmd(testOpts)
	cmd.SetArgs(append(testAccounts, "--wallet", tmpDir))
	err = cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	// Check that the accounts were deleted.
	for _, account := range testAccounts {
		_, err = wallet.GetAccount(account)
		if err == nil {
			t.Fatalf("account %s was not deleted", account)
		}
	}
}
