// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

// TaskRequest holds options for running a jaas task
type TaskRequest struct {
	Image       string
	Networks    []string
	Constraints []string
	EnvVars     []string
	Mounts      []string

	ShowLogs      bool
	Timeout       string
	RemoveService bool
	RegistryAuth  string
	Command       string
}
