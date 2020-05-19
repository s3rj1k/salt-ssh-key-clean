package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	debug *log.Logger
)

func init() {
	debug = log.New(
		ioutil.Discard,
		"# ",
		0,
	)

	if strings.EqualFold(os.Getenv("DEBUG"), "TRUE") {
		debug.SetOutput(os.Stderr)
	}
}
