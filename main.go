package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	cmdRosterFilePath     string
	cmdKnownHostsFilePath string
)

func init() {
	if os.Getuid() != 0 {
		fmt.Fprintf(os.Stderr, "%s needs to be run as root!\n", os.Args[0])

		os.Exit(1)
	}
}

func main() {
	flag.StringVar(&cmdRosterFilePath, "roster", "/root/roster_salt", "define an alternative location for the default roster file")
	flag.StringVar(&cmdKnownHostsFilePath, "hosts", "/root/.ssh/known_hosts", "define an alternative location for the default known_hosts file")
	flag.Parse()

	roster, err := parseRoster(cmdRosterFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)

		os.Exit(1)
	}

	for k, v := range roster {
		debug.Printf("%s: %s\n", k, v.String())

		for _, el := range getKnownHostsRecord(v.Host, v.Port) {
			fmt.Fprintf(os.Stdout, "%s\n", el.String())
		}
	}
}
