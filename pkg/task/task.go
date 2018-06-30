package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	jtypes "github.com/alexellis/jaas/pkg/types"
)

func List() ([]jtypes.Task, error) {
	var c *client.Client
	var err error
	c, err = client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("is the Docker Daemon running? Error: %s", err.Error())
	}
	jFilters := filters.NewArgs()
	jFilters.Add("label", "jaas=true")

	options := types.ServiceListOptions{
		Filters: jFilters,
	}

	services, err := c.ServiceList(context.Background(), options)

	if err != nil {
		return nil, err
	}

	retTasks := []jtypes.Task{}
	for _, service := range services {

		taskFilter := filters.NewArgs()
		taskFilter.Add("_up-to-date", "true")
		// taskFilter.Add("label", "jaas=true")
		taskFilter.Add("service", service.Spec.Name)

		// taskFilter.Add("desired-state", "running")

		taskListOptions := types.TaskListOptions{
			Filters: taskFilter,
		}

		taskList, inspectErr := c.TaskList(context.Background(), taskListOptions)
		if err != nil {
			return nil, inspectErr
		}
		for _, taskItem := range taskList {
			retTasks = append(retTasks, jtypes.Task{
				Name:     service.Spec.Name + "#" + taskItem.ID,
				Replicas: *service.Spec.Mode.Replicated.Replicas,
				Status:   string(taskItem.Status.State),
			})
		}
	}

	return retTasks, nil
}

func Create(image, command string) (*jtypes.TaskCreateStatus, error) {

	var c *client.Client
	var err error
	c, err = client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("is the Docker Daemon running? Error: %s", err.Error())
	}

	spec := makeSpec(image, []string{})
	createOptions := types.ServiceCreateOptions{}

	if len(command) > 0 {
		spec.TaskTemplate.ContainerSpec.Command = strings.Split(command, " ")
	}

	createResponse, _ := c.ServiceCreate(context.Background(), spec, createOptions)
	opts := types.ServiceInspectOptions{InsertDefaults: true}

	service, _, _ := c.ServiceInspectWithRaw(context.Background(), createResponse.ID, opts)
	fmt.Printf("Service created: %s (%s)\n", service.Spec.Name, createResponse.ID)

	// taskExitCode := pollTask(c, createResponse.ID, timeoutVal, taskRequest.ShowLogs, taskRequest.RemoveService)

	return &jtypes.TaskCreateStatus{ID: createResponse.ID}, nil
}

func makeSpec(image string, envVars []string) swarm.ServiceSpec {
	max := uint64(1)

	labels := map[string]string{
		"jaas": "true",
	}

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
		Annotations: swarm.Annotations{
			Labels: labels,
		},
	}
	return spec
}

func readEnvs(file string) ([]string, error) {
	var err error
	var envs []string

	data, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		return envs, readErr
	}

	lines := strings.Split(string(data), "\n")
	for n, line := range lines {
		if len(line) > 0 {
			if strings.Index(line, "=") == -1 {
				err = fmt.Errorf("no seperator found in line %d of env-file %s", n, file)
				break
			}
			envs = append(envs, line)
		}
	}
	return envs, err
}

const swarmError = -999
const timeoutError = -998
