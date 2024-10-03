package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

func worker(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		time.Sleep(500 * time.Millisecond)
		reversed := reverseString(task)
		fmt.Printf("Worker %d обработал: %s -> %s\n", id, task, reversed)
		results <- reversed
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	var numWorkers int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&numWorkers)

	file, err := os.Open("6.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	tasks := make(chan string, 10)
	results := make(chan string, 10)

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			tasks <- line
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	outputFile, err := os.Create("6out.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	for result := range results {
		outputFile.WriteString(result + "\n")
	}

	fmt.Println("Все задачи обработаны.")
}
