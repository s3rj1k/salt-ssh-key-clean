package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	ping "github.com/sparrc/go-ping"
)

const (
	timeout        = 30
	defaultSSHPort = 22
)

var (
	sshPath           string
	sshKeyScanBinPath string
)

type knownHost struct {
	Host   string
	Type   string
	PubKey string
}

func (s knownHost) String() string {
	return fmt.Sprintf("%s %s %s", s.Host, s.Type, s.PubKey)
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

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func testTCPPing(host string, port int) bool {
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(host, strconv.Itoa(port)),
		timeout*time.Second,
	)
	if err != nil {
		return false
	}

	defer conn.Close()

	return conn != nil
}

func testICMPPing(host string) bool {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return false
	}

	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Interval = timeout * time.Second
	pinger.Timeout = timeout * time.Second

	pinger.Run()

	stats := pinger.Statistics()

	return stats.PacketLoss == 0
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
