package main

import (
	"bufio"
	"io"

	"strings"

	"github.com/davecgh/go-spew/spew"
)

// ssh-keyscan -t rsa,dsa,ecdsa,ed25519 -p 2211 -H noc.mirohost.net

/*
root@salt:/srv/pillar/users# ssh-keygen -F [titan1.mirohost.net]:2211
root@salt:/srv/pillar/users# ssh-keygen -F [titan1.mirohost.net]:22
root@salt:/srv/pillar/users# ssh-keygen -F titan1.mirohost.net
# Host titan1.mirohost.net found: line 14536
|1|k0sOQHpR3q28P/opaWUUL1Dbxos=|88qaLm2ytZmTCDgvUDe5g4HuPo0= ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQRCWwLzT5HTbbsPkHu9lY2Wm9nx3R5GmbZi6CKQy7SPBFlxrZ7dGOeXgWQyetGg8ukXKXDPygZhtjkRGEEEwOkVMUSI0wQfN6q+nkn9m5+q4n3dOKW0D49uvoYX8WxIOZ9Gq3I2us3LZlcJawWIeQRLbcjVYRf2EBTrcZPUV7xQAIvZs8dq48r3oXnN71XCE55G8jVWm5RS3l5ljqHE7UmqBMvZVy4UPlBvObS91HeERB6jI0HGQ23ScDro2x1AYAN7x7Baqx2Y7u3iXNX8RXsjrD4pploO37RN2MimCa0hAToOUGTL9QdLqhjymMAQD0YpMO0qFyK7vXn2fWKxzN
*/

func main() {
	spew.Dump(sshKeyScan("noc.mirohost.net", 2211))
	spew.Dump(sshKeyFind("noc.mirohost.net", 2211))
}

type KnownHost struct {
	Host string
	Type string
	Key  string
}

func toKnownHosts(r io.Reader) chan KnownHost {
	out := make(chan KnownHost)

	go func(r io.Reader, out chan KnownHost) {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "# ") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) != 3 {
				continue
			}

			out <- KnownHost{
				Host: fields[0],
				Type: fields[1],
				Key:  fields[2],
			}
		}

		close(out)

		if err := scanner.Err(); err != nil {
			return
		}
	}(r, out)

	return out
}
