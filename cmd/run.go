// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	jswarm "github.com/alexellis/jaas/pkg/swarm"
	jtypes "github.com/alexellis/jaas/pkg/types"
	"github.com/spf13/cobra"
)

var (
	taskRequest jtypes.TaskRequest
	verbose     bool
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVarP(&taskRequest.Image, "image", "i", "", "Docker image to run")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.Networks, "network", "n", []string{}, "provide a network to bind to")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.Constraints, "constraint", "c", []string{}, "constraint for task")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.EnvVars, "env", "e", []string{}, "environmental variable for task")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.Mounts, "mount", "m", []string{}, "bind-mount a volume the task")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.Secrets, "secret", "s", []string{}, "Add existing secret to task")

	runCmd.PersistentFlags().StringArrayVar(&taskRequest.EnvFiles, "env-file", []string{}, "populate environment from an envfile for the task")

	runCmd.PersistentFlags().BoolVarP(&taskRequest.ShowLogs, "show-logs", "l", true, "show logs")
	runCmd.PersistentFlags().StringVarP(&taskRequest.Timeout, "timeout", "t", "60s", "timeout as a Golang duration")

	runCmd.PersistentFlags().StringVarP(&taskRequest.Command, "command", "d", "", "Command to run")

	runCmd.PersistentFlags().BoolVarP(&taskRequest.RemoveService, "remove", "r", true, "remove service after running task")
	runCmd.PersistentFlags().StringVarP(&taskRequest.RegistryAuth, "registry", "a", "", "registry auth string in base64")
	runCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "b", false, "print verbose debug information")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a task",
	Long:  `Run a task in Docker Swarm using a Service and task`,
	RunE:  runRun,
}

func runRun(cmd *cobra.Command, args []string) error {
	err := jswarm.RunTask(taskRequest)
	return err
}
