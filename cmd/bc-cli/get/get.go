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

package get

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/bestchains/bc-cli/pkg/account"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/connProfile"
	"github.com/bestchains/bc-cli/pkg/depository"
	"github.com/bestchains/bc-cli/pkg/federation"
	"github.com/bestchains/bc-cli/pkg/org"
	"github.com/bestchains/bc-cli/pkg/proposal"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
	}
	cmd.AddCommand(depository.NewGetDepositoryCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	cmd.AddCommand(account.NewGetAccountCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	cmd.AddCommand(org.NewOrgGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	cmd.AddCommand(connProfile.NewGetConnProfileCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	cmd.AddCommand(federation.NewFedGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	//cmd.AddCommand(chaincode.NewCCGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	//cmd.AddCommand(chaincodebuild.NewCCBGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	//cmd.AddCommand(channel.NewChanGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	//cmd.AddCommand(vote.NewVoteGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	cmd.AddCommand(proposal.NewProposalGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	//cmd.AddCommand(policy.NewPolicyGetCmd(common.Options{IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}}))
	return cmd
}
