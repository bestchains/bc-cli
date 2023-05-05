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
	"fmt"
	"os"
	"testing"

	"github.com/bestchains/bc-cli/pkg/common"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	delTestBasePath = ".delete"
	delTestPath     = getTestBasePath + "/wallet"
)

func TestNewDeleteAccountCmd(t *testing.T) {
	// step 1: create 4 accounts.
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	createCmd := NewCreateAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: b1, ErrOut: b2}})
	_ = createCmd.Flags().Set("wallet", getTestPath)
	for i := 0; i < 4; i++ {
		_ = createCmd.Execute()
	}

	deleteScanfFormat := "account \"%s\" deleted\n"
	errFormat := "Error: account \"%s\" remove %s/%s: no such file or directory\n"
	bufOutput := bytes.NewBuffer([]byte{})
	bufErrOutput := bytes.NewBuffer([]byte{})
	delCmd := NewDeleteAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: bufOutput, ErrOut: bufErrOutput}})
	_ = delCmd.Flags().Set("wallet", delTestPath)

	dirEntries, err := os.ReadDir(delTestPath)
	if err != nil {
		t.Fatalf("run read dir error %s", err)
	}
	expectOutput := bytes.NewBuffer([]byte{})
	output := bytes.NewBuffer([]byte{})

	normalDelFiles := make([]string, 0)
	for _, dir := range dirEntries {
		if dir.IsDir() {
			continue
		}
		expectOutput.WriteString(fmt.Sprintf(deleteScanfFormat, dir.Name()))
		normalDelFiles = append(normalDelFiles, dir.Name())
	}

	// step 1: First delete all accounts.
	delCmd.SetArgs(normalDelFiles)
	if err := delCmd.Execute(); err != nil {
		t.Fatalf("run delete account cmd with args %v error %s", normalDelFiles, err)
	}
	output.Write(bufOutput.Bytes())

	// step 2: Delete non-existent account information
	missingAccount := []string{"abc", "def"}
	for _, account := range missingAccount {
		expectOutput.WriteString(fmt.Sprintf(errFormat, account, delTestPath, account))
	}

	delCmd.SetArgs(missingAccount)
	if err := delCmd.Execute(); err != nil {
		t.Fatalf("run delete account cmd with args %v error %s", missingAccount, err)
	}
	output.Write(bufErrOutput.Bytes())
	if expectOutput.String() != output.String() {
		t.Fatalf("expect \r\n%s get \r\n%s", expectOutput.String(), output.String())
	}

	os.RemoveAll(delTestBasePath)
}
