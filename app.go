package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	"io/ioutil"

	"golang.org/x/net/context"
)

// Enable Docker's experimental mode in 1.13-rc before continuing.
func main() {
	var image string
	var timeout int
	var showlogs bool

	flag.StringVar(&image, "image", "", "Docker image name")
	flag.BoolVar(&showlogs, "showlogs", true, "show logs from stdout")
	flag.IntVar(&timeout, "timeout", 0, "ticks until we time out the service")
	flag.Parse()

	if len(image) == 0 {
		fmt.Println("Provide an image name with -image")
		return
	}

	var c *client.Client
	var err error
	c, err = client.NewEnvClient()
	if err != nil {
		log.Fatal("Error with Docker client.")
	}

	spec := makeSpec(image)
	createOptions := types.ServiceCreateOptions{}
	createResponse, _ := c.ServiceCreate(context.Background(), spec, createOptions)
	fmt.Printf("Service created: %s\n", createResponse.ID)

	pollTask(c, createResponse.ID, timeout, showlogs)
}

func makeSpec(image string) swarm.ServiceSpec {
	max := uint64(1)
	spec := swarm.ServiceSpec{
		TaskTemplate: swarm.TaskSpec{
			RestartPolicy: &swarm.RestartPolicy{
				MaxAttempts: &max,
				Condition:   swarm.RestartPolicyConditionNone,
			},
			ContainerSpec: swarm.ContainerSpec{
				Image: image,
			},
		},
	}
	return spec
}

func pollTask(c *client.Client, id string, timeout int, showlogs bool) {
	filters2 := filters.NewArgs()
	filters2.Add("id", id)

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
			if showTasks(c, item.ID, showlogs) {
				return
			}

			if timeout > 0 && ticks >= timeout {
				fmt.Printf("Timing out after %d ticks.", ticks)
				return
			}
		}
	}
}

func showTasks(c *client.Client, id string, showLogs bool) bool {
	filters1 := filters.NewArgs()
	filters1.Add("service", id)
	// fmt.Println("Task filter: ", id)
	val, _ := c.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters1,
	})

	var done bool
	for _, v := range val {
		if v.Status.State == swarm.TaskStateComplete {
			fmt.Println("\n")
			fmt.Println("Printing service logs")
			fmt.Printf("Exit code: %d\n", v.Status.ContainerStatus.ExitCode)
			fmt.Printf("State: %d\n", v.Status.State)

			if showLogs == true {
				logRequest, _ := c.ServiceLogs(context.Background(), id, types.ContainerLogsOptions{
					Follow:     false,
					ShowStdout: true,
					ShowStderr: true,
					Timestamps: true,
					Details:    false,
				})
				defer logRequest.Close()

				//	, ShowStderr: true, ShowStdout: true})
				res, _ := ioutil.ReadAll(logRequest)

				fmt.Println(string(res[:]))
			}

			done = true
			break
		} else {
			fmt.Printf(".")
		}
	}
	return done
}
