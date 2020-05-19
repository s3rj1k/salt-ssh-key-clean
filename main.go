package main

import (
	"github.com/davecgh/go-spew/spew"
)

func main() {
	scanedKeys := sshKeyScan("noc.mirohost.net", 2211)
	availiableKeys := sshKeyFind("noc.mirohost.net", 2211)

	spew.Dump(
		intersectKnownHosts(scanedKeys, availiableKeys),
	)
}
