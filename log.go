package main

import (
	"io/ioutil"
	"log"
	"os"
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

	// if strings.EqualFold(os.Getenv("DEBUG"), "TRUE") {
	debug.SetOutput(os.Stderr)
	// }
}
