package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	cmdRosterFilePath        string
	cmdPrivateKeyPath        string
	cmdNumberOfConcurentJobs int
)

func main() {
	if os.Getuid() != 0 {
		fatal.Fatalf("%s needs to be run as root!\n", os.Args[0])
	}

	flag.StringVar(&cmdRosterFilePath, "roster", "/etc/salt/roster", "defines location for the default roster file")
	flag.StringVar(&cmdPrivateKeyPath, "privkey", "/etc/salt/pki/master/ssh/salt-ssh.rsa", "defines location for the default private key")
	flag.IntVar(&cmdNumberOfConcurentJobs, "parallel", 5, "defines amount of concurent workers")

	flag.Parse()

	targets, err := getTargetsFromRoster(cmdRosterFilePath)
	if err != nil {
		fatal.Fatalf("%s\n", err)
	}

	for target := range worker(
		cmdNumberOfConcurentJobs,
		targets,
	) {
		fmt.Printf("%s\n", target)
	}
}
