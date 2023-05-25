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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"github.com/bestchains/bestchains-contracts/library"
	"github.com/bestchains/bestchains-contracts/library/context"
)

// Account represents a user account with an address and private key
type Account struct {
	Address    string `json:"address"`
	PrivateKey []byte `json:"privKey"`

	// signer parsed from PrivateKey
	signer ecdsa.PrivateKey
}

// NewAccount generates a new account with a new public-private key pair
// Returns a new account and an error if the key generation fails
func NewAccount() (*Account, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Generate a new address from the public key
	addr := new(library.Address)
	if err = addr.FromPublicKey(&pk.PublicKey); err != nil {
		return nil, err
	}

	// Convert the private key to an encoded PEM block
	x509Encoded, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return nil, err
	}

	// Return a new account with the generated address and private key
	return &Account{
		Address:    addr.String(),
		PrivateKey: pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}),
	}, nil
}

// GenerateAndSignMessage generates a context.Message and signs it with the account's signer.
// args is a variadic parameter that can take multiple strings.
// Returns the base64-encoded string representation of the message and an error, if any.
func (account *Account) GenerateAndSignMessage(nonce uint64, args ...string) (string, error) {
	// create a new message with the given nonce
	msg := context.Message{
		Nonce:     nonce,
		PublicKey: "",
		Signature: "",
	}

	// generate the signature for the message with the account's signer and the given args
	err := msg.GenerateSignature(&account.signer, args...)
	if err != nil {
		// return an empty string and the error if there was an error generating the signature
		return "", err
	}

	// return the base64-encoded string representation of the message
	return msg.Base64EncodedStr()
}
