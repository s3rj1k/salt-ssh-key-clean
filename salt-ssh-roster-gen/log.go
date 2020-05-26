package main

import (
	"log"
	"os"
)

var (
	// debug *log.Logger
	fatal *log.Logger
)

func init() {
	// debug = log.New(
	// 	ioutil.Discard,
	// 	"DBG:",
	// 	log.LstdFlags|log.Lshortfile,
	// )

	fatal = log.New(
		os.Stderr,
		"FTL:",
		log.LstdFlags,
	)

	// if strings.EqualFold(os.Getenv("DEBUG"), "TRUE") {
	// 	debug.SetOutput(os.Stderr)
	// }
}
