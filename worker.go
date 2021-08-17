package main

import (
	"sync"
)

func worker(maxWorkers int, in <-chan target) <-chan target {
	const (
		unavailableByTCPandICMP = "[-] (TCP, ICMP)"
		unavailableByTCP        = "[-] (TCP)"
		unknownKey              = "[-] (SSH)"
		valid                   = "[+]"
	)

	wg := sync.WaitGroup{}

	out := make(chan target, cap(in))
	defer close(out)

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)

		go func(in <-chan target, out chan<- target, worker int) {
			defer func() {
				wg.Done()
			}()

			for target := range in {
				if !testTCPPing(target.Host, target.Port) {
					if !testICMPPing(target.Host) {
						// out <- target // we are going to remove only explicety invalid keys
						critical.Printf("%s #[%03d]: %v\n", unavailableByTCPandICMP, worker, target)

						continue
					}

					// out <- target // we are going to remove only explicety invalid keys
					critical.Printf("%s #[%03d]: %v\n", unavailableByTCP, worker, target)

					continue
				}

				if !testSSHKey(cmdPrivateKeyPath, target.Host, target.User, target.Port) {
					out <- target
					critical.Printf("%s #[%03d]: %v\n", unknownKey, worker, target)

					continue
				}

				info.Printf("%s #[%03d]: %v\n", valid, worker, target)
			}
		}(in, out, w)
	}

	wg.Wait()

	return out
}
