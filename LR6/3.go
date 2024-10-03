package main

import (
	"fmt"
	"math/rand"
	"time"
)

func genereateRandomNumbers(n int, chInt chan<- int) {
	fmt.Println("Генерация случайных чисел...")
	for i := 0; i < n; i++ {
		time.Sleep(500 * time.Millisecond)
		num := rand.Intn(100)
		chInt <- num
	}
	close(chInt)
	fmt.Println("Генерация случайных чисел завершена!")
}

func isEven(chInt <-chan int, chString chan<- string) {
	for v := range chInt {
		if (v % 2 == 0) {
			chString <- fmt.Sprintf("%d четное", v)
		} else {
			chString <- fmt.Sprintf("%d нечетное", v)
		}
	}
	close(chString)
}

func main() {
	chInt := make(chan int)
	chString := make(chan string)
	
	go genereateRandomNumbers(10, chInt)
	go isEven(chInt, chString)
	
	for {
		select {
		case number, ok := <-chInt:
			if ok {
				fmt.Printf("Сгенерировано число: %d\n", number)
			} else {
				chInt = nil
			}
		case result, ok := <-chString:
			if ok {
				fmt.Println(result)
			} else {
				chString = nil
			}
		}
		if chInt == nil && chString == nil {
			break
		}
	}
}