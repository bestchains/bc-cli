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

package main

import (
	"github.com/bestchains/bc-cli/cmd/bc-cli/create"
	delcmd "github.com/bestchains/bc-cli/cmd/bc-cli/delete"
	"github.com/bestchains/bc-cli/cmd/bc-cli/get"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bc-cli [usage]",
		Short: "Command line tools for Bestchains",
		Long:  `bc-cli is a command tool for bestchain that can query and create a variety of blockchain resources.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(create.NewCreateCmd())
	cmd.AddCommand(get.NewGetCmd())
	cmd.AddCommand(delcmd.NewDeleteCmd())
	return cmd
}

func main() {
	if err := NewCmd().Execute(); err != nil {
		panic(err)
	}
}
