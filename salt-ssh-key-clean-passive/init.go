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
		"[DBG]",
		0,
	)

	critical = log.New(
		os.Stderr,
		"[CRT]",
		0,
	)

	fatal = log.New(
		os.Stderr,
		"[FTL]",
		0,
	)

	if strings.EqualFold(os.Getenv("DEBUG"), "TRUE") {
		debug.SetOutput(os.Stderr)
	}

	var err error

	if sshKeyScanBinPath, err = exec.LookPath("ssh-keyscan"); err != nil {
		fatal.Fatalf("ssh-keyscan binary not found!\n")
	}

	if sshKeyGenBinPath, err = exec.LookPath("ssh-keygen"); err != nil {
		fatal.Fatalf("ssh-keygen binary not found!\n")
	}
}