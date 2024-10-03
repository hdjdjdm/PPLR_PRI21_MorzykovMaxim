package main

import (
	"fmt"
	"sync"
)

var counter int = 0

func main() {
	n := 10
	wg := sync.WaitGroup{}
	wg.Add(n)
	var mutex sync.Mutex

	for i := 1; i <= n; i++ {
		go work(i, &mutex, &wg)
	}
	
	wg.Wait()
	fmt.Println("\nКонец")
}

func work(num int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mutex.Lock()
	counter++
	fmt.Printf("Горутина %d: %d\n", num, counter)
	mutex.Unlock()
}