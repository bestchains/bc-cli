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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	// Create a temporary directory for the wallet
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Create a new local wallet
	wallet, err := NewLocalWallet(tempDir)
	assert.NoError(t, err)

	// Generate a new account
	account, err := NewAccount()
	assert.NoError(t, err)

	// Store the account in the wallet
	err = wallet.StoreAccount(account)
	assert.NoError(t, err)

	// Retrieve the account from the wallet
	loadedAccount, err := wallet.GetAccount(account.Address)
	assert.NoError(t, err)
	assert.Equal(t, account.Address, loadedAccount.Address)
	assert.Equal(t, account.PrivateKey, loadedAccount.PrivateKey)

	// List the accounts in the wallet
	accounts, err := wallet.ListAccounts()
	assert.NoError(t, err)
	assert.Equal(t, []string{account.Address}, accounts)

	// Delete the account from the wallet
	filePath := tempDir + "/" + account.Address
	err = wallet.DeleteAccounts(account.Address)
	assert.NoError(t, err)
	assert.NoFileExists(t, filePath)
}
