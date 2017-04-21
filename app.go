package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
		log.Fatal("Is the Docker Daemon running?")
		return
	}

	// Check that experimental mode is enabled on the daemon, fall back to no logging if not
	if showlogs {
		if versionInfo, _ := c.ServerVersion(context.Background()); !versionInfo.Experimental {
			fmt.Println("Experimental daemon required to display service logs, falling back to no log display.")
			showlogs = false
		}
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
	opts := types.ServiceInspectOptions{InsertDefaults: true}

	service, _, _ := c.ServiceInspectWithRaw(context.Background(), createResponse.ID, opts)
	fmt.Printf("Service created: %s (%s)\n", service.Spec.Name, createResponse.ID)

	taskExitCode := pollTask(c, createResponse.ID, timeout, showlogs, removeService)
	os.Exit(taskExitCode)

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
	return exitCode, done
}
