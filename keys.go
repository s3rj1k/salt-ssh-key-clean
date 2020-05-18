package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

const timeout = 5

var (
	sshKeyScanBinPath string
	sshKeyGenBinPath  string
)

func init() {
	var err error

	if sshKeyScanBinPath, err = exec.LookPath("ssh-keyscan"); err != nil {
		log.Fatal(err)
	}

	if sshKeyGenBinPath, err = exec.LookPath("ssh-keygen"); err != nil {
		log.Fatal(err)
	}
}

func knownHostExecOutputWrapper(name string, args ...string) []KnownHost {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		name,
		args...,
	)

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}

	cmd.Start()

	out := make([]KnownHost, 0)

	for el := range toKnownHosts(pipe) {
		out = append(out, el)
	}

	cmd.Wait()

	return out
}

func sshKeyFind(host string, port int) []KnownHost {
	var search string

	if port > 0 {
		search = fmt.Sprintf("[%s]:%d", host, port)
	} else {
		search = host
	}

	return knownHostExecOutputWrapper(sshKeyGenBinPath, "-F", search)
}

func sshKeyScan(host string, port int) []KnownHost {
	args := make([]string, 0)
	args = append(
		args,
		"-t", "rsa,dsa,ecdsa,ed25519",
	)

	if port > 0 {
		args = append(
			args,
			"-p", strconv.Itoa(port),
		)
	}

	args = append(
		args,
		host,
	)

	return knownHostExecOutputWrapper(sshKeyScanBinPath, args...)
}
