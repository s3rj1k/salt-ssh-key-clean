package main

import (
	"sync"
)

func worker(maxWorkers int, in <-chan target) <-chan []knownHost {
	wg := sync.WaitGroup{}

	out := make(chan []knownHost, len(in))
	defer close(out)

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)

		go func(in <-chan target, out chan<- []knownHost, worker int) {
			defer func() {
				wg.Done()
			}()

			for target := range in {
				if !testTCPPing(target.Host, target.Port) {
					if !testICMPPing(target.Host) {
						info.Printf("Unavailable host (TCP, ICMP): %v\n", target)
						debug.Printf("[-] {%d} %v\n", worker, target)

						continue
					}

					info.Printf("Unavailable host (TCP): %v\n", target)
					debug.Printf("[-] {%d} %v\n", worker, target)

					continue
				}

				if !testSSHKey(cmdPrivateKeyPath, target.Host, target.User, target.Port) {
					info.Printf("Unknown host key: %v\n", target)
					debug.Printf("[-] {%d} %v\n", worker, target)

					continue
				}

				debug.Printf("[+] {%d} %v\n", worker, target)
				out <- sshKeyScan(target.Host, target.Port)
			}
		}(in, out, w)
	}

	wg.Wait()

	return out
}
