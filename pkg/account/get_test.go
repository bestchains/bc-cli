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
	"reflect"
	"testing"

	"github.com/bestchains/bc-cli/pkg/common"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	getTestBasePath = ".get"
	getTestPath     = getTestBasePath + "/wallet"
)

func TestNewGetAccountCmd(t *testing.T) {
	// step 1: create 4 accounts.
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	createCmd := NewCreateAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: b1, ErrOut: b2}})
	_ = createCmd.Flags().Set("wallet", getTestPath)
	for i := 0; i < 4; i++ {
		_ = createCmd.Execute()
	}

	bufOutput := bytes.NewBuffer([]byte{})
	bufErrOutput := bytes.NewBuffer([]byte{})
	getCmd := NewGetAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: bufOutput, ErrOut: bufErrOutput}})
	expectOutput := []string{fmt.Sprintf("Error: stat %s: no such file or directory\nError: stat /tmp/def/abc: no such file or directory\n", common.WalletConfigDir)}

	dirEntries, err := os.ReadDir(getTestPath)
	if err != nil {
		t.Fatalf("run read dir error %s", err)
	}

	buf := bytes.NewBuffer([]byte("ACCOUNT\n"))
	for _, dir := range dirEntries {
		if dir.IsDir() {
			continue
		}
		buf.WriteString(fmt.Sprint(dir.Name(), "\n"))
	}
	expectOutput = append(expectOutput, buf.String())

	output := make([]string, 0)
	// step 1: Use the default path and the default path does not exist
	if err := getCmd.Execute(); err != nil {
		t.Fatalf("run get account cmd with default wallet %s error %s", common.WalletConfigDir, err)
	}
	// step 2: Using non-existent paths to obtain account information
	_ = getCmd.Flags().Set("wallet", "/tmp/def/abc")
	if err := getCmd.Execute(); err != nil {
		t.Fatalf("run get account cmd with wallet /tmp/def/abc error %s", err)
	}
	output = append(output, bufErrOutput.String())

	// step 3: Use the correct path to get the correct account information output
	_ = getCmd.Flags().Set("wallet", getTestPath)
	// Create a new directory and file, the get command should ignore it.
	_ = os.Mkdir(getTestPath+"/xyz", 0755)
	_, _ = os.Create(getTestPath + "/xyz/zyx")

	if err := getCmd.Execute(); err != nil {
		t.Fatalf("run get account cmd with wallet %s error %s", getTestPath, err)
	}
	output = append(output, bufOutput.String())
	if !reflect.DeepEqual(expectOutput, output) {
		t.Fatalf("expect %v get %v", expectOutput, output)
	}

	os.RemoveAll(getTestBasePath)
}
