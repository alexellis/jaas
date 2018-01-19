// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Execute run the root cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "jaas",
	Short: "Run an ad-hoc command",
	Long: `
Run ad-hoc commands`,
	RunE: executeRoot,
}

func executeRoot(cmd *cobra.Command, args []string) error {
	fmt.Printf(figletStr)
	cmd.Help()

	return nil
}

func init() {

}
