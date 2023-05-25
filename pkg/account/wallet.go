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
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type IWallet interface {
	StoreAccount(*Account) error
	GetAccount(string) (*Account, error)
	ListAccounts() ([]string, error)
	DeleteAccounts(...string) error
}

var _ IWallet = (*LocalWallet)(nil)

type LocalWallet struct {
	home string
}

// NewLocalWallet creates a new LocalWallet instance with the given home directory.
// If the directory does not exist, it will be created.
// The function returns a LocalWallet instance and an error if the directory creation fails.
func NewLocalWallet(home string) (LocalWallet, error) {
	home = strings.TrimSuffix(home, "/")             // remove trailing slash if present
	if _, err := os.Stat(home); os.IsNotExist(err) { // check if directory exists
		err = os.MkdirAll(home, os.ModePerm) // create directory with all permissions
		if err != nil {
			return LocalWallet{}, errors.Wrap(err, "mkdir local wallet home dir") // return error with context
		}
	}

	return LocalWallet{home: home}, nil // return LocalWallet instance and nil error
}

// StoreAccount stores the account information in a file with the address as the filename
func (localWallet *LocalWallet) StoreAccount(account *Account) error {
	// Convert the account to a JSON byte slice
	bytes, err := json.Marshal(account)
	if err != nil {
		return errors.Wrap(err, "invalid account")
	}

	// Create the target file path
	targetFile := filepath.Join(localWallet.home, account.Address)

	// Open the target file for writing
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open target account file")
	}
	defer file.Close()

	// Write the account bytes to the file
	_, err = file.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to write target account file")
	}

	return nil
}

// GetAccount retrieves an account by its address.
func (localWallet *LocalWallet) GetAccount(accAddr string) (*Account, error) {
	// Construct the file path
	filePath := filepath.Join(localWallet.home, accAddr)

	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find account file")
	}

	// Read account info from file
	objBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read account file")
	}

	var account = new(Account)

	// Unmarshal JSON data into Account object
	err = json.Unmarshal(objBytes, &account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account file")
	}

	// Check that the account address matches the given account name
	if account.Address != accAddr {
		return nil, errors.Errorf("expected account %s but got %s", account, account.Address)
	}

	// Parse private key
	pkEncoded, _ := pem.Decode(account.PrivateKey)
	pk, err := x509.ParseECPrivateKey(pkEncoded.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse account private key")
	}
	account.signer = *pk

	return account, nil
}

// ListAccounts returns a slice of account addresses stored in the local wallet directory.
// Each account address is represented as a string.
// An error is returned if the directory cannot be read.
func (localWallet *LocalWallet) ListAccounts() ([]string, error) {
	// Read the directory entries in the wallet directory.
	dirEntries, err := os.ReadDir(localWallet.home)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read wallet dir")
	}

	// Create an empty slice to store account addresses.
	accountAddrs := make([]string, 0, len(dirEntries))

	// Iterate over the directory entries.
	for _, info := range dirEntries {
		// Skip directories.
		if info.IsDir() {
			continue
		}
		// Append the file name (account address) to the account addresses slice.
		accountAddrs = append(accountAddrs, info.Name())
	}

	// Return the account addresses slice and nil error.
	return accountAddrs, nil
}

// DeleteAccounts deletes the accounts with the given addresses from the local wallet.
func (localWallet *LocalWallet) DeleteAccounts(accAddrs ...string) error {
	for _, accAddr := range accAddrs {
		filePath := filepath.Join(localWallet.home, accAddr)
		if err := os.Remove(filePath); err != nil {
			return errors.Wrapf(err, "failed to delete account %s", accAddr)
		}
	}
	return nil
}
