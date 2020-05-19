package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	cmdRosterFilePath     string
	cmdKnownHostsFilePath string
)

func init() {
	if os.Getuid() != 0 {
		log.Fatalf("%s needs to be run as root!\n", os.Args[0])
	}
}

func main() {
	flag.StringVar(&cmdRosterFilePath, "roster", "/root/roster_salt", "define an alternative location for the default roster file")
	flag.StringVar(&cmdKnownHostsFilePath, "hosts", "/root/.ssh/known_hosts", "define an alternative location for the default known_hosts file")
	flag.Parse()

	roster, err := parseRoster(cmdRosterFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range roster {
		fmt.Printf("# %s: %s\n", k, v.String())
		fmt.Printf("%s\n", getKnownHostsRecord(v.Host, v.Port))
	}
}
