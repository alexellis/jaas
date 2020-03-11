// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"os"
)

func main() {
	hostname, _ := os.Hostname()
	fmt.Println("Hostname: " + hostname)
}
