package main

import (
	"sync"
)

func worker(maxWorkers int, targets <-chan target, knownHostsPath string) []knownHost {
	wg := sync.WaitGroup{}

	hosts := make(chan knownHost)
	out := make([]knownHost, 0)

	go func(hosts <-chan knownHost) {
		for host := range hosts {
			out = append(out, host)
		}
	}(hosts)

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)

		go func(targets <-chan target, hosts chan<- knownHost, knownHostsPath string) {
			defer func() {
				wg.Done()
			}()

			for target := range targets {
				for _, el := range getKnownHostsRecord(target.Host, target.Port, knownHostsPath) {
					hosts <- el
				}
			}
		}(targets, hosts, knownHostsPath)
	}

	wg.Wait()
	close(hosts)

	return out
}
