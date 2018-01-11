package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

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
	var registry string
	var constraints listFlag
	var envVars listFlag

	flag.Var(&envVars, "env", "environmental variables")

	flag.StringVar(&image, "image", "", "Docker image name")
	flag.StringVar(&network, "network", "", "Docker swarm network name")

	flag.BoolVar(&showlogs, "showlogs", true, "show logs from stdout")
	flag.BoolVar(&removeService, "rm", false, "remove service after completion")
	flag.IntVar(&timeout, "timeout", 60, "ticks until we time out the service - default is 60 seconds")

	flag.StringVar(&registry, "registryAuth", "", "pass your registry authentication")

	flag.Var(&constraints, "constraint", "Placement constraints (e.g. node.labels.key==value)")

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
	versionInfo, versionErr := c.ServerVersion(context.Background())
	if versionErr != nil {
		log.Fatal("Is the Docker Daemon running?")

		return
	}

	if showlogs {
		apiVersion, parseErr := strconv.ParseFloat(versionInfo.APIVersion, 32)
		if parseErr != nil {
			apiVersion = 0
		}
		if apiVersion < 1.29 && versionInfo.Experimental == false {
			fmt.Println("Experimental daemon or Docker API Version 1.29+ required to display service logs, falling back to no log display.")
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

	if len(registry) > 0 {
		createOptions.EncodedRegistryAuth = registry
		fmt.Println("Using auth: " + registry)
	}

	placement := &swarm.Placement{}
	if len(constraints) > 0 {
		placement.Constraints = constraints
		spec.TaskTemplate.Placement = placement
	}

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
			ContainerSpec: &swarm.ContainerSpec{
				Image: image,
				Env:   *envVars,
			},
		},
	}
	return spec
}
