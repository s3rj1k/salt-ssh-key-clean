package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	cmdRosterFilePath        string
	cmdKnownHostsFilePath    string
	cmdNumberOfConcurentJobs int
)

func main() {
	if os.Getuid() != 0 {
		fatal.Fatalf("%s needs to be run as root!\n", os.Args[0])
	}

	flag.StringVar(&cmdRosterFilePath, "roster", "/etc/salt/roster", "defines an location for the default roster file")
	flag.StringVar(&cmdKnownHostsFilePath, "hosts", "/root/.ssh/known_hosts", "defines an location for the default known_hosts file")
	flag.IntVar(&cmdNumberOfConcurentJobs, "parallel", 5, "defines amount of concurent workers")

	flag.Parse()

	targets, err := getTargetsFromRoster(cmdRosterFilePath)
	if err != nil {
		fatal.Fatalf("%s\n", err)
	}

	f, err := ioutil.TempFile(
		filepath.Dir(cmdKnownHostsFilePath),
		filepath.Base(cmdKnownHostsFilePath)+".",
	)
	if err != nil {
		fatal.Fatalf("known_hosts: %s\n", err)
	}

	defer func(f *os.File) {
		if err := f.Sync(); err != nil {
			critical.Printf("known_hosts: %s\n", err)
		}
		if err := f.Close(); err != nil {
			fatal.Fatalf("known_hosts: %s\n", err)
		}
	}(f)

	hosts := make(chan knownHost)

	go func(hosts <-chan knownHost, f *os.File) {
		for host := range hosts {
			if _, err := f.WriteString(host.StringLn()); err != nil {
				fatal.Fatalf("%s\n", err)
			}
		}
	}(hosts, f)

	worker(
		cmdNumberOfConcurentJobs,
		targets,
		hosts,
		cmdKnownHostsFilePath,
	)
}
