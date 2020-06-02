package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	timeout        = 5
	defaultSSHPort = 22
)

type knownHost struct {
	Host   string
	Type   string
	PubKey string
}

func (s knownHost) String() string {
	return fmt.Sprintf("%s %s %s", s.Host, s.Type, s.PubKey)
}

func (s knownHost) StringLn() string {
	return s.String() + "\n"
}

func (s knownHost) KeyWithType() string {
	return fmt.Sprintf("%s %s", s.Type, s.PubKey)
}

var (
	sshKeyScanBinPath string
	sshKeyGenBinPath  string
)

func knownHostExecOutputWrapper(name string, args ...string) []knownHost {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		name,
		args...,
	)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		debug.Printf("exec: %v\n", err)

		return nil
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		debug.Printf("exec: %v\n", err)

		return nil
	}

	if err := cmd.Start(); err != nil {
		debug.Printf("exec: %v\n", err)

		return nil
	}

	out := make([]knownHost, 0)

	for el := range toKnownHosts(stdoutPipe, stderrPipe) {
		out = append(out, el)
	}

	if err := cmd.Wait(); err != nil {
		debug.Printf("exec: %v\n", err)

		return nil
	}

	return out
}

func sshKeyFind(host string, port int, knownHostsPath string) []knownHost {
	var search string

	if port > 0 && port != defaultSSHPort {
		search = fmt.Sprintf("[%s]:%d", host, port)
	} else {
		search = host
	}

	return knownHostExecOutputWrapper(
		sshKeyGenBinPath,
		"-F", search,
		"-f", knownHostsPath,
	)
}

func sshKeyScan(host string, port int, _ string) []knownHost {
	args := make([]string, 0)
	args = append(
		args,
		"-t", "rsa,dsa,ecdsa,ed25519",
	)

	if port > 0 && port != defaultSSHPort {
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

func toKnownHosts(readers ...io.Reader) <-chan knownHost {
	out := make(chan knownHost)
	r := io.MultiReader(readers...)

	go func(r io.Reader, out chan knownHost) {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "@") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 3 || len(fields) > 4 {
				continue
			}

			out <- knownHost{
				Host:   fields[0],
				Type:   fields[1],
				PubKey: fields[2],
			}
		}

		close(out)

		if err := scanner.Err(); err != nil {
			debug.Printf("scanner: %v\n", err)

			return
		}
	}(r, out)

	return out
}

func intersectKnownHosts(left, right []knownHost) []knownHost {
	if len(left) == 0 || len(right) == 0 {
		return nil
	}

	intersected := make([]knownHost, 0, len(left)+len(right))

	for i := range left {
		for j := range right {
			if left[i].KeyWithType() == right[j].KeyWithType() {
				intersected = append(intersected, left[i])
			}
		}
	}

	return intersected
}

func getKnownHostsRecord(host string, port int, knownHostsPath string) []knownHost {
	return intersectKnownHosts(
		sshKeyScan(host, port, knownHostsPath),
		sshKeyFind(host, port, knownHostsPath),
	)
}
