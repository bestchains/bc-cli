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
	"testing"

	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/stretchr/testify/assert"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestNewCreateAccountCmd(t *testing.T) {
	// Create a temporary wallet directory
	tempDir := t.TempDir()

	// Create a buffer for capturing command output
	output := new(bytes.Buffer)

	// Create a dummy options struct
	options := common.Options{
		IOStreams: genericclioptions.IOStreams{
			Out:    output,
			ErrOut: output,
		},
	}

	// Create the Cobra command
	cmd := NewCreateAccountCmd(options)
	cmd.SetArgs([]string{"--wallet", tempDir})

	// Execute the command
	err := cmd.Execute()

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the account file was created in the wallet directory
	wallet, err := NewLocalWallet(tempDir)
	assert.NoError(t, err)

	accs, err := wallet.ListAccounts()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(accs))

	// Assert that the command output contains the expected message

	expectedOutput := fmt.Sprintf("account/%s created\n", accs[0])
	assert.Equal(t, expectedOutput, output.String())
}
