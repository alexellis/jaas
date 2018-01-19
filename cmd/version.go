// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"runtime"

	"github.com/alexellis/jaas/version"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

// GitCommit injected at build-time
var (
	shortVersion bool
)

func init() {
	versionCmd.Flags().BoolVar(&shortVersion, "short-version", false, "Just print Git SHA")

	rootCmd.AddCommand(versionCmd)
}

// versionCmd displays version information
var versionCmd = &cobra.Command{
	Use:   "version [--short-version]",
	Short: "Display the clients version information",
	Long: fmt.Sprintf(`The version command returns the current clients version information.

This currently consists of the GitSHA from which the client was built.
- https://github.com/alexellis/jaas/tree/%s`, version.GitCommit),
	Example: `  jaas version
  jaas version --short-version`,
	Run: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	if shortVersion {
		fmt.Println(version.BuildVersion())
	} else {
		figletColoured := aec.BlueF.Apply(figletStr)
		if runtime.GOOS == "windows" {
			figletColoured = aec.GreenF.Apply(figletStr)
		}
		fmt.Printf(figletColoured)
		fmt.Printf("Commit: %s\n", version.GitCommit)
		fmt.Printf("Version: %s\n", version.BuildVersion())
	}
}

const figletStr = `
     _             ____  
    | | __ _  __ _/ ___| 
 _  | |/ _` + "`" + ` |/ _` + "`" + `  \___ \ 
| |_| | (_| | (_| |___) |
 \___/ \__,_|\__,_|____/ 
 
`
