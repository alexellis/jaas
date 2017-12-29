package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func showTasks(c *client.Client, id string, showLogs, removeService bool) (int, bool) {
	filters1 := filters.NewArgs()
	filters1.Add("service", id)
	// fmt.Println("Task filter: ", id)
	val, _ := c.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters1,
	})

	exitCode := 1
	var done bool
	for _, v := range val {
		if v.Status.State == swarm.TaskStateComplete {
			fmt.Println("\n")
			fmt.Printf("Exit code: %d\n", v.Status.ContainerStatus.ExitCode)
			fmt.Printf("State: %s\n", v.Status.State)
			fmt.Println("\n")

			exitCode = v.Status.ContainerStatus.ExitCode
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
					fmt.Printf("Unable to pull service logs.\nError: %s", err)
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
