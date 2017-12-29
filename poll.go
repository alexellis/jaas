package main

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const swarmError = -999
const timeoutError = -998

func pollTask(c *client.Client, id string, timeout int, showlogs, removeService bool) int {
	filters2 := filters.NewArgs()
	filters2.Add("id", id)

	exitCode := swarmError

	opts := types.ServiceListOptions{
		Filters: filters2,
	}

	list, _ := c.ServiceList(context.Background(), opts)
	for _, item := range list {
		ticks := 0
		fmt.Println("ID: ", item.ID, " Update at: ", item.UpdatedAt)
		for {
			time.Sleep(500 * time.Millisecond)
			ticks++
			taskExitCode, found := showTasks(c, item.ID, showlogs, removeService)
			if found {
				exitCode = taskExitCode
				break
			}

			if timeout > 0 && ticks >= timeout {
				fmt.Printf("Timing out after %d ticks.", ticks)
				return timeoutError
			}
		}
	}

	return exitCode
}
