package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func monteCarloPiPart(points int, wg *sync.WaitGroup, resultChan chan<- int) {
	defer wg.Done()
	inside := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < points; i++ {
		x, y := r.Float64(), r.Float64()
		if x*x+y*y <= 1 {
			inside++
		}
	}
	resultChan <- inside
}

func runMonteCarlo(totalPoints int, nWorkers int) (float64, time.Duration) {
	pointsPerWorker := totalPoints / nWorkers
	var wg sync.WaitGroup
	resultChan := make(chan int, nWorkers)

	start := time.Now()

	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go monteCarloPiPart(pointsPerWorker, &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	totalInside := 0
	for val := range resultChan {
		totalInside += val
	}

	duration := time.Since(start)
	pi := 4.0 * float64(totalInside) / float64(pointsPerWorker*nWorkers)
	return pi, duration
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // використовуємо всі ядра

	N := 1_000_000
	threads := []int{1, 2, 4, 8, 16, 32, 64}

	fmt.Printf("Обчислення числа π методом Монте-Карло (%d точок)\n", N)
	for _, n := range threads {
		pi, dur := runMonteCarlo(N, n)
		fmt.Printf("%2d горутин: pi ≈ %.6f, час = %v\n", n, pi, dur)
	}
}
