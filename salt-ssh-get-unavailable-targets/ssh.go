package main

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

const (
	timeout        = 5
	defaultSSHPort = 22
)

var (
	sshPath string
)

func testPing(key, host, user string, port int) bool {
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

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
