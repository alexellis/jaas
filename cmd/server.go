// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/alexellis/jaas/pkg/proxy"
	"github.com/alexellis/jaas/pkg/server"
	"github.com/alexellis/jaas/pkg/types"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	username, address, image, command string
	passwdStdin                       bool
)

func init() {
	rootCmd.AddCommand(serverCmd)

	loginServerCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	loginServerCmd.PersistentFlags().BoolVarP(&passwdStdin, "password-stdin", "p", false, "password passed via stdin")
	loginServerCmd.PersistentFlags().StringVarP(&address, "address", "a", "http://127.0.0.1:8081", "remote server address")

	runServerCmd.PersistentFlags().StringVarP(&image, "image", "i", "", "Docker image")
	runServerCmd.PersistentFlags().StringVarP(&command, "command", "c", "", "Docker command override")

	serverCmd.AddCommand(listServerCmd)
	serverCmd.AddCommand(startServerCmd)
	serverCmd.AddCommand(stopServerCmd)
	serverCmd.AddCommand(loginServerCmd)
	serverCmd.AddCommand(runServerCmd)

}

var loginServerCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into jaas server",
	Long:  `Log into jaas server`,
	RunE:  loginServerRun,
}
var listServerCmd = &cobra.Command{
	Use:   "list",
	Short: "List jaas remote jobs",
	Long:  `List jaas remote jobs`,
	RunE:  listServerRun,
}
var runServerCmd = &cobra.Command{
	Use:   "run",
	Short: "Run jaas command remotely",
	Long:  `Run jaas command remotely`,
	RunE:  runServerRun,
}

var stopServerCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the jaas server",
	Long:  `Stop the jaas server`,
	RunE:  stopServerRun,
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the jaas server",
	Long:  `Start the jaas server`,
	RunE:  startServerRun,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server commands",
	Long:  `Commands for jaas as a server`,
}

func listServerRun(cmd *cobra.Command, args []string) error {
	var err error

	fmt.Printf("List")

	config, configErr := readConfig()
	if configErr != nil {
		return configErr
	}

	var remoteAddress string

	if len(address) > 0 {
		remoteAddress = address
	} else {
		return fmt.Errorf("no --address given")
	}

	client := proxy.NewJaaSClient(*config)
	items, err := client.List(remoteAddress)
	for _, item := range items {
		fmt.Println(item.Name, item.Replicas, item.Status)
	}

	return err
}

func runServerRun(cmd *cobra.Command, args []string) error {
	var err error
	if len(image) == 0 {
		return fmt.Errorf("expected --image flag")
	}

	fmt.Printf("Running: %s %s\n", image, command)

	config, configErr := readConfig()
	if configErr != nil {
		return configErr
	}

	var remoteAddress string

	if len(address) > 0 {
		remoteAddress = address
	} else {
		return fmt.Errorf("no --address given")
	}

	client := proxy.NewJaaSClient(*config)
	statusCode, err := client.Run(remoteAddress, proxy.RunRequest{
		Job: &proxy.JaasJob{
			Image:   image,
			Command: command,
		},
	})

	fmt.Printf("Status: %d\n", statusCode)

	return err
}

func stopServerRun(cmd *cobra.Command, args []string) error {
	var err error
	fmt.Println("Stopping")
	return err
}

func readConfig() (*types.JaaSConfig, error) {

	home := os.Getenv("HOME")
	dest := path.Join(home, ".jaas")
	filename := path.Join(dest, "credentials.yaml")

	bytesYaml, yErr := ioutil.ReadFile(filename)
	if yErr != nil {
		return nil, yErr
	}

	jaasConfig := types.JaaSConfig{}
	yamlErr := yaml.Unmarshal(bytesYaml, &jaasConfig)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return &jaasConfig, nil
}

func startServerRun(cmd *cobra.Command, args []string) error {
	var err error
	fmt.Println("Starting")

	jaasConfig, configErr := readConfig()
	if configErr != nil {
		return configErr
	}

	server := server.NewJaaSServer(8081, time.Second*10, *jaasConfig)
	stopCh := make(chan interface{})
	err = server.Start(stopCh)
	if err != nil {
		return err
	}

	return err
}

func loginServerRun(cmd *cobra.Command, args []string) error {
	var err error

	if len(username) == 0 {
		return fmt.Errorf("give a --username")
	}
	if passwdStdin == false {
		return fmt.Errorf("give a password via stdin with --password-stdin")
	}

	fmt.Printf("Enter password: ")
	passBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	mkErr := prepareDirectory()
	if mkErr != nil {
		return mkErr

	}

	home := os.Getenv("HOME")
	dest := path.Join(home, ".jaas")
	filename := path.Join(dest, "credentials.yaml")
	encodedPass := hex.EncodeToString(passBytes)

	config := types.JaaSConfig{
		Auths: []types.JaaSServerAuth{},
	}

	config.Auths = append(config.Auths, types.JaaSServerAuth{
		Address:  address,
		Username: username,
		Password: encodedPass,
	})

	yamlBytes, marshalErr := yaml.Marshal(config)
	if marshalErr != nil {
		return marshalErr
	}

	writeFileErr := ioutil.WriteFile(filename, yamlBytes, 0700)
	if writeFileErr != nil {
		return writeFileErr
	}

	fmt.Printf("Logging in as %s\n", username)

	client := proxy.NewJaaSClient(config)
	pingRes, err := client.Ping(address)

	if err != nil {
		return err
	}

	fmt.Printf("Login: %d\n", pingRes.Status)

	return err
}

func prepareDirectory() error {
	home := os.Getenv("HOME")
	dest := path.Join(home, ".jaas")
	if _, err := os.Stat(dest); err != nil {
		mkErr := os.MkdirAll(dest, 0700)
		if mkErr != nil {
			return mkErr
		}
	}
	return nil
}
