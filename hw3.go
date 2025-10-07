package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func collatzSteps(n uint64) uint64 {
	var steps uint64 = 0
	for n > 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	return steps
}

func worker(id int, jobs <-chan uint64, totalSteps *uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	for n := range jobs {
		steps := collatzSteps(n)
		atomic.AddUint64(totalSteps, steps)
	}
}

func main() {
	var (
		workers int
		maxN    uint64
	)
	flag.IntVar(&workers, "workers", 0, "number of worker goroutines (default: number of CPUs)")
	flag.Uint64Var(&maxN, "max", 10000000, "maximum natural number to compute (inclusive)")
	flag.Parse()

	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	fmt.Printf("Collatz parallel run: max=%d, workers=%d\n", maxN, workers)

	jobs := make(chan uint64, workers*4)
	var wg sync.WaitGroup
	var totalSteps uint64 = 0

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(i, jobs, &totalSteps, &wg)
	}

	for i := uint64(1); i <= maxN; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	avg := float64(totalSteps) / float64(maxN)
	fmt.Printf("Done. Total steps = %d\n", totalSteps)
	fmt.Printf("Average steps per number = %.6f\n", avg)
}
