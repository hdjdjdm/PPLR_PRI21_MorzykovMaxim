package main

import (
	"fmt"
	"sync"
	"time"
)

func fibonacci(n int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	a, b := 0, 1
	for i := 0; i < n; i++ {
		result <- a
		a, b = b, a+b
		time.Sleep(100 * time.Millisecond)
	}
	close(result)
	fmt.Println("\nГенерация чисел Фибоначчи завершена. Канал закрыт.")
}

func readChannel(result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range result {
		fmt.Printf("%d ", v)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	result := make(chan int)
	
	fmt.Println("\nГенерация чисел Фибоначчи...")
	go fibonacci(10, result, &wg)
	go readChannel(result, &wg)
	
	wg.Wait()
}