package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("The final value is not guaranteed to be 1000 because multiple goroutines perform an unsynchronized read-modify-write operation on the same shared variable, causing a data race and lost updates.")
	fmt.Println(counter)
}
