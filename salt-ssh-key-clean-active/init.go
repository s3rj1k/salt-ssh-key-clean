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

	critical = log.New(
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

	if sshKeyScanBinPath, err = exec.LookPath("ssh-keyscan"); err != nil {
		fatal.Fatalf("ssh-keyscan binary not found!\n")
	}
}
