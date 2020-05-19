package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
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

func (s knownHost) KeyWithType() string {
	return fmt.Sprintf("%s %s", s.Type, s.PubKey)
}

var (
	sshKeyScanBinPath string
	sshKeyGenBinPath  string
)

func init() {
	var err error

	if sshKeyScanBinPath, err = exec.LookPath("ssh-keyscan"); err != nil {
		fmt.Fprintf(os.Stderr, "ssh-keyscan binary not found!\n")

		os.Exit(1)
	}

	if sshKeyGenBinPath, err = exec.LookPath("ssh-keygen"); err != nil {
		fmt.Fprintf(os.Stderr, "ssh-keygen binary not found!\n")

		os.Exit(1)
	}
}

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

	return deDuplicateKnownHosts(out)
}

func sshKeyFind(host string, port int) []knownHost {
	var search string

	if port > 0 && port != defaultSSHPort {
		search = fmt.Sprintf("[%s]:%d", host, port)
	} else {
		search = host
	}

	return knownHostExecOutputWrapper(
		sshKeyGenBinPath,
		"-F", search,
		"-f", cmdKnownHostsFilePath,
	)
}

func sshKeyScan(host string, port int) []knownHost {
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

// https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable
func deDuplicateKnownHosts(s []knownHost) []knownHost {
	if len(s) == 0 {
		return nil
	}

	sort.Slice(
		s, func(i, j int) bool {
			return s[i].KeyWithType() < s[j].KeyWithType()
		},
	)

	j := 0

	for i := 1; i < len(s); i++ {
		if s[i].KeyWithType() == s[j].KeyWithType() {
			continue
		}

		j++

		s[j] = s[i]
	}

	return s[:j+1]
}

func intersectKnownHosts(left, right []knownHost) []knownHost {
	if len(left) == 0 || len(right) == 0 {
		return nil
	}

	intersected := make([]knownHost, 0, len(left)+len(right))

outer:
	for i := range deDuplicateKnownHosts(left) {
		for j := range deDuplicateKnownHosts(right) {
			if left[i].KeyWithType() == right[j].KeyWithType() {
				intersected = append(intersected, left[i])

				continue outer
			}
		}
	}

	return intersected
}

func getKnownHostsRecord(host string, port int) []knownHost {
	knownHosts := intersectKnownHosts(
		sshKeyScan(host, port),
		sshKeyFind(host, port),
	)

	return knownHosts
}
