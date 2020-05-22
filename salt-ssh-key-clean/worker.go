package main

import (
	"fmt"
	"os"
	"sync"
)

func worker(maxWorkers int, targets <-chan target, knownHostsPath string) {
	wg := sync.WaitGroup{}

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)

		go func(targets <-chan target, knownHostsPath string) {
			defer func() {
				wg.Done()
			}()

			for target := range targets {
				for _, el := range getKnownHostsRecord(target.Host, target.Port, knownHostsPath) {
					fmt.Fprintf(os.Stdout, "%s\n", el.String())
				}
			}
		}(targets, knownHostsPath)
	}

	wg.Wait()
}
