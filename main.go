package main

import (
	"flag"
	"os"
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

	const (
		removeKeyError = "[-] (KEY)"
		removeKeyOK    = "[+] (KEY)"
	)

	flag.StringVar(&cmdRosterFilePath, "roster", "/etc/salt/roster", "defines location for the default roster file")
	flag.StringVar(&cmdKnownHostsFilePath, "hosts", "/root/.ssh/known_hosts", "defines an location for the default known_hosts file")
	flag.StringVar(&cmdPrivateKeyPath, "privkey", "/etc/salt/pki/master/ssh/salt-ssh.rsa", "defines location for the default private key")
	flag.IntVar(&cmdNumberOfConcurentJobs, "parallel", 25, "defines amount of concurent workers")

	flag.Parse()

	targets, err := getTargetsFromRoster(cmdRosterFilePath)
	if err != nil {
		fatal.Fatalf("%s\n", err)
	}

	sshKeygenRemoveArgs := make([]string, 0, len(targets))

	for t := range worker(
		cmdNumberOfConcurentJobs,
		targets,
	) {
		args := t.RemoveList()
		if len(args) > 0 {
			sshKeygenRemoveArgs = append(sshKeygenRemoveArgs, args...)
		}
	}

	for _, el := range sshKeygenRemoveArgs {
		if err := removeSSHKey(cmdKnownHostsFilePath, el); err != nil {
			critical.Printf("%s: %s\n", removeKeyError, el)
		} else {
			info.Printf("%s: %s\n", removeKeyOK, el)
		}
	}
}
