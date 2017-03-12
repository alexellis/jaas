package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	"io/ioutil"

	"golang.org/x/net/context"
)

type listFlag []string

func (i *listFlag) String() string {
	str := ""
	for _, v := range *i {
		str += v
	}
	return str
}

func (i *listFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Enable Docker's experimental mode in 1.13-rc before continuing.
func main() {
	var image string
	var timeout int
	var showlogs bool
	var network string
	var removeService bool

	var envVars listFlag

	flag.Var(&envVars, "env", "environmental variables")

	flag.StringVar(&image, "image", "", "Docker image name")
	flag.StringVar(&network, "network", "", "Docker swarm network name")

	flag.BoolVar(&showlogs, "showlogs", true, "show logs from stdout")
	flag.BoolVar(&removeService, "rm", false, "remove service after completion")
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

	spec := makeSpec(image, &envVars)
	if len(network) > 0 {
		nets := []swarm.NetworkAttachmentConfig{
			swarm.NetworkAttachmentConfig{Target: network},
		}
		spec.Networks = nets
	}

	createOptions := types.ServiceCreateOptions{}
	createResponse, _ := c.ServiceCreate(context.Background(), spec, createOptions)

	service, _, _ := c.ServiceInspectWithRaw(context.Background(), createResponse.ID)
	fmt.Printf("Service created: %s (%s)\n", service.Spec.Name, createResponse.ID)

	pollTask(c, createResponse.ID, timeout, showlogs, removeService)
}

func makeSpec(image string, envVars *listFlag) swarm.ServiceSpec {
	max := uint64(1)

	spec := swarm.ServiceSpec{
		TaskTemplate: swarm.TaskSpec{
			RestartPolicy: &swarm.RestartPolicy{
				MaxAttempts: &max,
				Condition:   swarm.RestartPolicyConditionNone,
			},
			ContainerSpec: swarm.ContainerSpec{
				Image: image,
				Env:   *envVars,
			},
		},
	}
	return spec
}

func pollTask(c *client.Client, id string, timeout int, showlogs, removeService bool) {
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
			if showTasks(c, item.ID, showlogs, removeService) {
				return
			}

			if timeout > 0 && ticks >= timeout {
				fmt.Printf("Timing out after %d ticks.", ticks)
				return
			}
		}
	}
}

func showTasks(c *client.Client, id string, showLogs, removeService bool) bool {
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
			fmt.Printf("Exit code: %d\n", v.Status.ContainerStatus.ExitCode)
			fmt.Printf("State: %s\n", v.Status.State)
			fmt.Println("\n")

			if showLogs {
				fmt.Println("Printing service logs")
			}

			if showLogs {
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
	return done
}
