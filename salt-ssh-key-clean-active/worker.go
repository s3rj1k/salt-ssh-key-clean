package main

import (
	"runtime"
	"sync"
)

func worker(maxWorkers int, in <-chan target) <-chan []knownHost {
	wg := sync.WaitGroup{}

	out := make(chan []knownHost, len(in))
	defer close(out)

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)

		go func(in <-chan target, out chan<- []knownHost) {
			defer func() {
				wg.Done()
			}()

			for target := range in {
				if testActivePing(cmdPrivateKeyPath, target.Host, target.User, target.Port) {
					debug.Printf("[+] {%d} %v\n", runtime.NumGoroutine(), target)

					out <- sshKeyScan(target.Host, target.Port)
				} else {
					info.Printf("Unavailable host: %v\n", target)
					debug.Printf("[-] {%d} %v\n", runtime.NumGoroutine(), target)
				}
			}
		}(in, out)
	}

	wg.Wait()

	return out
}
