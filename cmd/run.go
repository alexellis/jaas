// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	"github.com/spf13/cobra"
)

var (
	taskRequest TaskRequest
	verbose     bool
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVarP(&taskRequest.Image, "image", "i", "", "Docker image to run")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.Networks, "network", "n", []string{}, "provide a network to bind to")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.Constraints, "constraint", "c", []string{}, "constraint for task")
	runCmd.PersistentFlags().StringArrayVarP(&taskRequest.EnvVars, "env", "e", []string{}, "environmental variable for task")
	runCmd.PersistentFlags().BoolVarP(&taskRequest.ShowLogs, "show-logs", "l", true, "show logs")
	runCmd.PersistentFlags().StringVarP(&taskRequest.Timeout, "timeout", "t", "60s", "timeout as a Golang duration")

	runCmd.PersistentFlags().BoolVarP(&taskRequest.RemoveService, "remove", "r", true, "remove service after running task")
	runCmd.PersistentFlags().StringVarP(&taskRequest.RegistryAuth, "registry", "a", "", "registry auth string in base64")
	runCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print verbose debug information")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a task",
	Long:  `Run a task in Docker Swarm using a Service and task`,
	RunE:  runRun,
}

func runRun(cmd *cobra.Command, args []string) error {
	err := runTask(taskRequest)
	return err
}

func validate(taskRequest TaskRequest) error {
	if len(taskRequest.Image) == 0 {
		return fmt.Errorf("must a valid supply --image")
	}
	return nil
}

func runTask(taskRequest TaskRequest) error {
	if validationErr := validate(taskRequest); validationErr != nil {
		return validationErr
	}

	if verbose {
		fmt.Printf("Running.. OK %s\n", taskRequest.Image)
		fmt.Printf("Connected to.. OK %s\n", taskRequest.Networks)
		fmt.Printf("Constraints: %s\n", taskRequest.Constraints)
		fmt.Printf("envVars: %s\n", taskRequest.EnvVars)
	}

	timeoutVal, parseErr := time.ParseDuration(taskRequest.Timeout)
	if parseErr != nil {
		return parseErr
	}

	if verbose {
		fmt.Printf("timeout: %s\n", timeoutVal)
	}

	var c *client.Client
	var err error
	c, err = client.NewEnvClient()
	if err != nil {

		return fmt.Errorf("is the Docker Daemon running? Error: %s", err.Error())
	}

	// Check that experimental mode is enabled on the daemon, fall back to no logging if not
	versionInfo, versionErr := c.ServerVersion(context.Background())
	if versionErr != nil {
		log.Fatal("Is the Docker Daemon running?")

		return versionErr
	}

	if taskRequest.ShowLogs {
		apiVersion, parseErr := strconv.ParseFloat(versionInfo.APIVersion, 32)
		if parseErr != nil {
			apiVersion = 0
		}
		if apiVersion < 1.29 && versionInfo.Experimental == false {
			return fmt.Errorf("experimental daemon or Docker API Version 1.29+ required to display service logs, falling back to no log display")
		}
	}

	spec := makeSpec(taskRequest.Image, taskRequest.EnvVars)
	if len(taskRequest.Networks) > 0 {
		nets := []swarm.NetworkAttachmentConfig{
			swarm.NetworkAttachmentConfig{Target: taskRequest.Networks[0]},
		}
		spec.Networks = nets
	}

	createOptions := types.ServiceCreateOptions{}

	if len(taskRequest.RegistryAuth) > 0 {
		createOptions.EncodedRegistryAuth = taskRequest.RegistryAuth
		fmt.Println("Using RegistryAuth")
	}

	placement := &swarm.Placement{}
	if len(taskRequest.Constraints) > 0 {
		placement.Constraints = taskRequest.Constraints
		spec.TaskTemplate.Placement = placement
	}

	createResponse, _ := c.ServiceCreate(context.Background(), spec, createOptions)
	opts := types.ServiceInspectOptions{InsertDefaults: true}

	service, _, _ := c.ServiceInspectWithRaw(context.Background(), createResponse.ID, opts)
	fmt.Printf("Service created: %s (%s)\n", service.Spec.Name, createResponse.ID)

	taskExitCode := pollTask(c, createResponse.ID, timeoutVal, taskRequest.ShowLogs, taskRequest.RemoveService)
	os.Exit(taskExitCode)
	return nil
}

func makeSpec(image string, envVars []string) swarm.ServiceSpec {
	max := uint64(1)

	spec := swarm.ServiceSpec{
		TaskTemplate: swarm.TaskSpec{
			RestartPolicy: &swarm.RestartPolicy{
				MaxAttempts: &max,
				Condition:   swarm.RestartPolicyConditionNone,
			},
			ContainerSpec: &swarm.ContainerSpec{
				Image: image,
				Env:   envVars,
			},
		},
	}
	return spec
}

const swarmError = -999
const timeoutError = -998

func pollTask(c *client.Client, id string, timeout time.Duration, showlogs, removeService bool) int {
	svcFilters := filters.NewArgs()
	svcFilters.Add("id", id)

	exitCode := swarmError

	opts := types.ServiceListOptions{
		Filters: svcFilters,
	}

	list, _ := c.ServiceList(context.Background(), opts)
	for _, item := range list {
		start := time.Now()
		end := start.Add(timeout)

		fmt.Println("ID: ", item.ID, " Update at: ", item.UpdatedAt)
		for {
			time.Sleep(500 * time.Millisecond)

			taskExitCode, found := showTasks(c, item.ID, showlogs, removeService)
			if found {
				exitCode = taskExitCode
				break
			}
			now := time.Now()
			if now.After(end) {
				fmt.Printf("Timing out after %s.", timeout.String())
				return timeoutError
			}
		}
	}

	return exitCode
}

func showTasks(c *client.Client, id string, showLogs, removeService bool) (int, bool) {
	filters1 := filters.NewArgs()
	filters1.Add("service", id)

	tasks, _ := c.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters1,
	})

	exitCode := 1
	var done bool
	stopStates := []swarm.TaskState{
		swarm.TaskStateComplete,
		swarm.TaskStateFailed,
		swarm.TaskStateRejected,
	}

	for _, task := range tasks {

		stop := false
		for _, stopState := range stopStates {
			if task.Status.State == stopState {
				stop = true
				break
			}
		}

		if stop {
			fmt.Printf("\n\n")
			fmt.Printf("Exit code: %d\n", task.Status.ContainerStatus.ExitCode)
			fmt.Printf("State: %s\n", task.Status.State)
			fmt.Printf("\n\n")

			exitCode = task.Status.ContainerStatus.ExitCode

			if exitCode == 0 && task.Status.State == swarm.TaskStateRejected {
				exitCode = 255 // force non-zero exit for task rejected
			}

			if showLogs {
				fmt.Println("Printing service logs")
			}

			if showLogs {
				logRequest, err := c.ServiceLogs(context.Background(), id, types.ContainerLogsOptions{
					Follow:     false,
					ShowStdout: true,
					ShowStderr: true,
					Timestamps: true,
					Details:    false,
					Tail:       "all",
				})

				if err != nil {
					fmt.Printf("Unable to pull service logs.\nError: %s\n", err)
				} else {
					defer logRequest.Close()

					//	, ShowStderr: true, ShowStdout: true})
					res, _ := ioutil.ReadAll(logRequest)

					fmt.Println(string(res[:]))
				}
			}

			if removeService {
				fmt.Println("Removing service...")
				if err := c.ServiceRemove(context.Background(), id); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			done = true
			break
		} else {
			fmt.Printf(".")
		}
	}
	return exitCode, done
}
