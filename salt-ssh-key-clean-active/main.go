package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	cmdRosterFilePath        string
	cmdKnownHostsFilePath    string
	cmdPrivateKeyPath        string
	cmdNumberOfConcurentJobs int
)

func main() {
	if os.Getuid() != 0 {
		fatal.Fatalf("%s needs to be run as root!\n", os.Args[0])
	}

	flag.StringVar(&cmdRosterFilePath, "roster", "/etc/salt/roster", "defines location for the default roster file")
	flag.StringVar(&cmdKnownHostsFilePath, "hosts", "/root/.ssh/known_hosts", "defines an location for the default known_hosts file")
	flag.StringVar(&cmdPrivateKeyPath, "privkey", "/etc/salt/pki/master/ssh/salt-ssh.rsa", "defines location for the default private key")
	flag.IntVar(&cmdNumberOfConcurentJobs, "parallel", 25, "defines amount of concurent workers")

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
		if f != nil {
			if f.Fd() != uintptr(^uint(0)) { // checks that Fd is still valid (0xffffffffffffffff), https://man7.org/linux/man-pages/man2/eventfd.2.html
				if err := f.Sync(); err != nil {
					critical.Printf("known_hosts: %s\n", err)
				}
				if err := f.Close(); err != nil {
					fatal.Fatalf("known_hosts: %s\n", err)
				}
			}
		}
	}(f)

	for host := range worker(
		cmdNumberOfConcurentJobs,
		targets,
	) {
		for i := range host {
			if _, err := fmt.Fprintf(f, "%s\n", host[i]); err != nil {
				fatal.Fatalf("known_hosts: %s\n", err)
			}
		}
	}

	if err := f.Sync(); err != nil {
		critical.Printf("known_hosts: %s\n", err)
	}

	newKnownHostsPath := f.Name()

	if err := f.Close(); err != nil {
		fatal.Fatalf("known_hosts: %s\n", err)
	}

	if err := os.Rename(cmdKnownHostsFilePath, cmdKnownHostsFilePath+".old"); err != nil {
		fatal.Fatalf("known_hosts: %s\n", err)
	}

	if err := os.Rename(newKnownHostsPath, cmdKnownHostsFilePath); err != nil {
		fatal.Fatalf("known_hosts: %s\n", err)
	}
}
