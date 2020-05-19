package main

import (
	"github.com/davecgh/go-spew/spew"
)

func main() {
	s := sshKeyScan("noc.mirohost.net", 2211)

	s = append(s, s...)

	spew.Dump(sshKeyFind("noc.mirohost.net", 2211))

	s = deDuplicateKnownHosts(s)

	spew.Dump(s)
}
