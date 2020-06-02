package main

import (
	"sync"
)

func worker(maxWorkers int, targets <-chan target, hosts chan<- knownHost, knownHostsPath string) {
	wg := sync.WaitGroup{}

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
}
