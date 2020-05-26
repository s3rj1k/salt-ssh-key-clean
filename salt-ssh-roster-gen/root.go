package main

import (
	"fmt"
	"os"
)

func init() {
	if os.Getuid() != 0 {
		fmt.Fprintf(os.Stderr, "%s needs to be run as root!\n", os.Args[0])

		os.Exit(1)
	}
}
