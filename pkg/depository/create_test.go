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

package depository

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bestchains/bc-cli/pkg/account"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bestchains-contracts/library/context"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"testing"
)

const (
	getTestBasePath = ".get"
	getTestPath     = getTestBasePath + "/wallet"
)

func TestValueDepotGen(t *testing.T) {
	expectValDepot := ValueDepository{
		Name:             "test name",
		ContentType:      "test type",
		ContentID:        "test ID",
		TrustedTimestamp: "123456789",
		Platform:         "test platform",
	}

	genBase := generateValueDepotBase64("test name", "test type", "test ID", "test platform")

	decodeRes, err := base64.StdEncoding.DecodeString(genBase)
	if err != nil {
		t.Fatalf("decode generated base64 failed: " + err.Error())
	}

	var resValDepot ValueDepository
	err = json.Unmarshal(decodeRes, &resValDepot)
	if err != nil {
		t.Fatalf("unmarshal valDep from decoded base64 failed: " + err.Error())
	}

	if resValDepot.Name != expectValDepot.Name {
		t.Fatalf("generated valDep don't match. expect name '%s', got '%s'", expectValDepot.Name, resValDepot.Name)
	}
}

func TestMsgGen(t *testing.T) {
	expectMsg := context.Message{
		Nonce:     0,
		PublicKey: "",
		Signature: "",
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key failed: " + err.Error())
	}
	genMsg := generateMessageBase64(0, *key, []byte("testValue"))
	decodeRes, err := base64.StdEncoding.DecodeString(genMsg)
	if err != nil {
		t.Fatalf("decode generated base64 failed: " + err.Error())
	}

	err = expectMsg.GenerateSignature(key, "testValue")
	if err != nil {
		t.Fatalf("expectMsg generate sign failed: " + err.Error())
	}

	var resMsg context.Message
	err = json.Unmarshal(decodeRes, &resMsg)
	if err != nil {
		t.Fatalf("unmarshal valDep from decoded base64 failed: " + err.Error())
	}

	if resMsg.Nonce != expectMsg.Nonce {
		t.Fatalf("generated Msg not match, expect nonce '%v', got '%v'", expectMsg.Nonce, resMsg.Nonce)
	}
}

func TestGetNonce(t *testing.T) {

}

func TestGetWalletInfo(t *testing.T) {
	// step 1: create 4 accounts.
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	createCmd := account.NewCreateAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: b1, ErrOut: b2}})
	_ = createCmd.Flags().Set("wallet", getTestPath)
	for i := 0; i < 4; i++ {
		_ = createCmd.Execute()
	}

	dirEntries, err := os.ReadDir(getTestPath)
	if err != nil {
		t.Fatalf("run read dir error %s", err)
	}

	// this results as reading the last account created
	var acc string
	for _, dir := range dirEntries {
		if dir.IsDir() {
			continue
		}
		acc = dir.Name()
	}

	// step 1: Using non-existent paths to obtain account information
	_, err = getWalletInfo("/tmp/def/abc", acc)
	if err == nil {
		t.Fatalf("run getWalletInfo with %s /tmp/def/abc, expected err, no err returned", acc)
	}

	// step 2: Use the correct path to get the correct account information output
	obj, err := getWalletInfo(getTestPath, acc)
	if err != nil {
		t.Fatalf("run getWalletInfo with account %s, path %s, error: %s", acc, getTestPath, err)
	}

	if acc != obj.Address {
		t.Fatalf("expect %v get %v", acc, obj.Address)
	}

	os.RemoveAll(getTestBasePath)
}
