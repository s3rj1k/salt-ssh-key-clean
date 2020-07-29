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
	info     *log.Logger
	critical *log.Logger
	fatal    *log.Logger
)

func init() {
	debug = log.New(
		ioutil.Discard,
		"[DBG] ",
		log.LstdFlags,
	)

	info = log.New(
		os.Stdout,
		"[NFO] ",
		log.LstdFlags,
	)

	critical = log.New(
		os.Stderr,
		"[CRT] ",
		log.LstdFlags,
	)

	fatal = log.New(
		os.Stderr,
		"[FTL] ",
		log.LstdFlags,
	)

	if strings.EqualFold(os.Getenv("DEBUG"), "TRUE") {
		debug.SetOutput(os.Stderr)
	}

	var err error

	if sshPath, err = exec.LookPath("ssh"); err != nil {
		fatal.Fatalf("ssh binary not found!\n")
	}

	if sshKeygenPath, err = exec.LookPath("ssh-keygen"); err != nil {
		fatal.Fatalf("ssh-keygen binary not found!\n")
	}
}
