package main

import (
	"fmt"
	"sync"
)

// Added a third parameter j to the goroutine inside the worker function to prevent race conditions
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		jValue := j // Create a new variable to use in goroutine
		go func(j int) {
			var result int
			switch j % 3 {
			case 0:
				result = j // Removed the unused j=j*1 operation
			case 1:
				result = j * 2
			case 2:
				result = j * 3
			}
			results <- result
		}(jValue)
	}
}

func main() {
	jobs := make(chan int)
	results := make(chan int)

	var wg sync.WaitGroup // Wait group to ensure all jobs are dispatched before closing the channel

	// Added the variable i as a parameter to the goroutine to prevent race conditions
	for i := 1; i <= 1000000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				i += 99
			}
			jobs <- i
		}(i)
	}

	// Now we wait for all jobs to be dispatched, then close the channel
	go func() {
		wg.Wait()
		close(jobs)
	}()

	jobs2 := []int{}
	for w := 1; w < 1000; w++ {
		jobs2 = append(jobs2, w)
	}

	var wgWorkers sync.WaitGroup // Wait group to ensure all workers have finished before closing the results channel

	for _, w := range jobs2 {
		wgWorkers.Add(1)
		go func(w int) {
			worker(w, jobs, results)
			wgWorkers.Done()
		}(w)
	}

	// Close the results channel after all workers have finished their jobs
	go func() {
		wgWorkers.Wait()
		close(results)
	}()

	var sum int32 = 0
	for r := range results {
		sum += int32(r)
	}
	fmt.Println(sum)
}
