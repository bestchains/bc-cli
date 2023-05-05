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
	"sort"
	"strings"
	"testing"

	"github.com/bestchains/bc-cli/pkg/common"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	createTestBasePath = ".create"
	createTestPath     = createTestBasePath + "/wallet"
	userPrivateKey     = "./pk.pem"
	pkContent          = `-----BEGIN PRIVATE KEY-----
MHcCAQEEIDuaob5MQI3tl8H/Z8L+VIiKaER1r/aojZfeRapKpbBhoAoGCCqGSM49
AwEHoUQDQgAER6bI26M8/6cEwpHNm+wHq/wxU4ISG/2xfcyGeAsghx4hAUjVg9rr
XYwFcMEK3BTGtx7v6Ai2OhxK4wF6/jibOA==
-----END PRIVATE KEY-----`
)

func TestNewCreateAccountCmd(t *testing.T) {
	scanfFormat := "account/%s created"
	bufOutput := bytes.NewBuffer([]byte{})
	bufErrOutput := bytes.NewBuffer([]byte{})

	f, err := os.Create(userPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(pkContent))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// step 1: Automatic generation of three accounts
	createCmd := NewCreateAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: bufOutput, ErrOut: bufErrOutput}})
	_ = createCmd.Flags().Set("wallet", createTestPath)
	for i := 0; i < 3; i++ {
		if err := createCmd.Execute(); err != nil {
			t.Fatalf("run create account cmd error %s", err)
		}
	}

	// step 2: Create an account with an existing private key
	_ = createCmd.Flags().Set("pk", userPrivateKey)
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("run create account cmd with pk error %s", err)
	}

	output := strings.Split(strings.TrimSpace(bufOutput.String()), "\n")
	files := make([]string, len(output))
	for i, o := range output {
		var fileName string
		fmt.Sscanf(o, scanfFormat, &fileName)
		files[i] = fileName
	}
	sort.Strings(files)

	// step 3: Check for file matches
	dirEntries, err := os.ReadDir(createTestPath)
	if err != nil {
		t.Fatalf("run read dir error %s", err)
	}
	expectFiles := make([]string, 0)
	for _, dir := range dirEntries {
		if dir.IsDir() {
			continue
		}
		expectFiles = append(expectFiles, dir.Name())
	}
	if !reflect.DeepEqual(expectFiles, files) {
		t.Fatalf("expect %v get %v", expectFiles, files)
	}

	os.RemoveAll(createTestBasePath)
	os.Remove(userPrivateKey)
}
