package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	cmdRosterFilePath        string
	cmdKnownHostsFilePath    string
	cmdNumberOfConcurentJobs int
)

func init() {
	if os.Getuid() != 0 {
		fmt.Fprintf(os.Stderr, "%s needs to be run as root!\n", os.Args[0])

		os.Exit(1)
	}
}

func main() {
	flag.StringVar(&cmdRosterFilePath, "roster", "/etc/salt/roster", "defines an location for the default roster file")
	flag.StringVar(&cmdKnownHostsFilePath, "hosts", "/root/.ssh/known_hosts", "defines an location for the default known_hosts file")
	flag.IntVar(&cmdNumberOfConcurentJobs, "parallel", 5, "defines amount of concurent workers")

	flag.Parse()

	targets, err := getTargetsFromRoster(cmdRosterFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)

		os.Exit(1)
	}

	worker(
		cmdNumberOfConcurentJobs,
		targets,
		cmdKnownHostsFilePath,
	)
}
