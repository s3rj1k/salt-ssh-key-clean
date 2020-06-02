package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	debug    *log.Logger
	critical *log.Logger
	fatal    *log.Logger
)

func init() {
	debug = log.New(
		ioutil.Discard,
		"[DBG] ",
		0,
	)

	critical = log.New(
		os.Stderr,
		"[CRT] ",
		0,
	)

	fatal = log.New(
		os.Stderr,
		"[FTL] ",
		0,
	)

	if strings.EqualFold(os.Getenv("DEBUG"), "TRUE") {
		debug.SetOutput(os.Stderr)
	}

	var err error

	if sshPath, err = exec.LookPath("ssh"); err != nil {
		fatal.Fatalf("ssh binary not found!\n")
	}
}
