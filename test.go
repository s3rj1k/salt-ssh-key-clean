package main

import (
	"net"
	"strconv"
	"time"

	ping "github.com/sparrc/go-ping"
)

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
