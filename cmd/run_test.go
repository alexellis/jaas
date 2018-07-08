// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var validRequest = TaskRequest{
	Image:         "input-image",
	Networks:      []string{"net1", "net2"},
	Constraints:   []string{"node.id=2ivku8v2gvtg4", "engine.labels.operatingsystem==ubuntu 14.04"},
	EnvVars:       []string{"ev1=val1", "ev2=val2"},
	Mounts:        []string{"hostVol1=taskVol1", "hostVol2=taskVol2"},
	Secrets:       []string{"secret1", "secret2"},
	ShowLogs:      true,
	Timeout:       "12",
	RemoveService: true,
	RegistryAuth:  "true",
	Command:       "echo 'some output'",
}

type fakeClient struct {
	client.CommonAPIClient
	secrets []string
}

func (fk fakeClient) SecretList(ctx context.Context, sopt types.SecretListOptions) ([]swarm.Secret, error) {
	slist := []swarm.Secret{}
	for i, secret := range fk.secrets {
		a := swarm.Annotations{Name: secret}
		sspec := swarm.SecretSpec{Annotations: a}
		s := swarm.Secret{
			ID:   strconv.Itoa(i),
			Meta: swarm.Meta{},
			Spec: sspec,
		}

		slist = append(slist, s)
	}
	return slist, nil
}

func newClient(secrets []string) fakeClient {
	return fakeClient{secrets: secrets}
}

func contains(el string, array []string) bool {
	for _, e := range array {
		if el == e {
			return true
		}
	}
	return false
}

func TestMakeServiceSpecValid(t *testing.T) {
	c := newClient([]string{"secret1", "secret2", "secret3"})

	f1, _ := ioutil.TempFile("", "jaas_env")
	f1body := []byte("f1var1=val11\nf1var2=val12\n")
	f1.Write(f1body)
	f1.Sync()
	f2, _ := ioutil.TempFile("", "jaas_env")
	f2body := []byte("f2var1=val21\nf2var2=val22\n")
	f2.Write(f2body)
	f2.Sync()
	defer os.Remove(f1.Name())
	defer os.Remove(f2.Name())
	validRequest.EnvFiles = []string{f1.Name(), f2.Name()}

	spec := makeServiceSpec(validRequest, c)
	if spec.TaskTemplate.ContainerSpec.Image != "input-image" {
		t.Errorf("Container spec image should be %s, was %s", validRequest.Image, spec.TaskTemplate.ContainerSpec.Image)
	}

	// Test networks
	networkTargets := []string{}
	for _, n := range spec.Networks {
		networkTargets = append(networkTargets, n.Target)
	}
	if !reflect.DeepEqual(networkTargets, []string{"net1", "net2"}) {
		t.Errorf("Container spec networks should be %s, was %s", validRequest.Networks, networkTargets)
	}

	// Test env vars from input and env files
	inputEnv := []string{"ev1=val1", "ev2=val2", "f1var1=val11", "f1var2=val12", "f2var1=val21", "f2var2=val22"}
	for _, ev := range inputEnv {
		if !contains(ev, spec.TaskTemplate.ContainerSpec.Env) {
			t.Errorf("Container spec env should contain %s, but only has %s", ev, spec.TaskTemplate.ContainerSpec.Env)
		}
	}

	// Test constraints
	if !reflect.DeepEqual(spec.TaskTemplate.Placement.Constraints, []string{"node.id=2ivku8v2gvtg4", "engine.labels.operatingsystem==ubuntu 14.04"}) {
		t.Errorf("Container spec constraints should be %s, was %s", validRequest.Constraints, spec.TaskTemplate.Placement.Constraints)
	}

	// Test mounts
	expectedMounts := []mount.Mount{
		{Source: "hostVol1", Target: "taskVol1"},
		{Source: "hostVol2", Target: "taskVol2"},
	}
	if !reflect.DeepEqual(spec.TaskTemplate.ContainerSpec.Mounts, expectedMounts) {
		t.Error("Container spec mounts should include:")
		for _, m := range expectedMounts {
			t.Errorf("{Source: %s, Target: %s}", m.Source, m.Target)
		}
		t.Error("But contained instead:")
		for _, m := range spec.TaskTemplate.ContainerSpec.Mounts {
			t.Errorf("{Source: %s, Target: %s}", m.Source, m.Target)
		}
	}

	// Test secrets
	secretNames := []string{}
	for _, s := range spec.TaskTemplate.ContainerSpec.Secrets {
		secretNames = append(secretNames, s.SecretName)
	}
	if !reflect.DeepEqual(secretNames, []string{"secret1", "secret2"}) {
		t.Errorf("Container spec secrets should be %s, was %s", validRequest.Secrets, secretNames)
	}

	// Test command
	if !reflect.DeepEqual(spec.TaskTemplate.ContainerSpec.Command, []string{"echo", "'some", "output'"}) {
		t.Errorf("Container spec command should be %s, was %s", []string{"echo", "'some", "output'"}, spec.TaskTemplate.ContainerSpec.Command)
	}
}
