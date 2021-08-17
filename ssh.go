package main

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

const (
	timeout        = 30
	defaultSSHPort = 22
)

var (
	sshPath       string
	sshKeygenPath string
)

func testSSHKey(key, host, user string, port int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		sshPath,
		"-q",
		"-i",
		key,
		"-o",
		"StrictHostKeyChecking=yes",
		"-o",
		"PasswordAuthentication=no",
		"-o",
		fmt.Sprintf("ConnectTimeout=%d", timeout),
		"-l",
		user,
		"-p",
		strconv.Itoa(port),
		host,
		"-t",
		"/bin/true",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		debug.Printf("%s {%s %d}: %v, %v\n", sshPath, host, port, err, out)

		return false
	}

	return true
}

func removeSSHKey(key, host string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		sshKeygenPath,
		"-q",
		"-f",
		key,
		"-R",
		host,
	)

	if err := cmd.Run(); err != nil {
		debug.Printf("%s {%s}: %v\n", sshKeygenPath, host, err)

		return err
	}

	return nil
}
