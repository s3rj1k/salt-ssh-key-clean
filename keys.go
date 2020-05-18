package main

import (
	"bytes"
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

func sshKeyFind(host string, port int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	var search string

	if port > 0 {
		search = fmt.Sprintf("[%s]:%d", host, port)
	} else {
		search = host
	}

	out, err := exec.CommandContext(
		ctx,
		sshKeyGenBinPath,
		"-F", search,
	).CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(out, []byte("\n"))

	var i int

	for _, b := range lines {
		if !bytes.HasPrefix(b, []byte("# ")) {
			lines[i] = b
			i++
		}
	}

	lines = lines[:i]

	return bytes.Join(lines, []byte("\n")), nil
}

func sshKeyScan(host string, port int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

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

	cmd := exec.CommandContext(
		ctx,
		sshKeyScanBinPath,
		args...,
	)

	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Start()

	if err = toKnownHosts(output); err != nil {
		return nil, err
	}

	cmd.Wait()

	return nil, nil
}
