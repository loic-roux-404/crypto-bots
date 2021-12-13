/*
Copyright © 2021 loic-roux-404 loic.roux.404@gmail.com

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
package cmd

import (
	"log"

	"github.com/loic-roux-404/crypto-bots/cmd/core"
	"github.com/loic-roux-404/crypto-bots/pkg/strategy"

	"github.com/spf13/cobra"
)

// forwarderCmd represents the forwarder command
var forwarderCmd = &cobra.Command{
	Use:   "forwarder",
	Short: "Forward native currency to another wallet",
	Long: `Instantly forward native currency deposited in crypto wallet
 to another configured one`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Info: forwarder started...")
		strategy.Forwarder(args[len(args)-1])
	},
}

func init() {
	RootCmd.AddCommand(core.InitChainCmd(forwarderCmd))
}
