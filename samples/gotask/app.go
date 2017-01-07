package main

import (
	"fmt"
        "os"
)

func main() {
        hostname, _ := os.Hostname()
        fmt.Println("Hostname: " + hostname)
}
