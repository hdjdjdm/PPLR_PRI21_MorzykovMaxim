package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func factorial(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Факториал начинает вычисляться...")
	time.Sleep(1500 * time.Millisecond)
	
	if n < 0 {
		fmt.Println("Ошибка: факториал отрицательного числа не определен.")
		return
	}
	
	if n == 0 {
		return
	}
	
	result := 1
	
	for i := 2; i <= n; i++ {
		result *= i
	}
	
	fmt.Printf("Факториал 4 = %d\n", result)
}

func genereateRandomNumbers(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Генерация случайных чисел...")
	for i := 0; i < n; i++ {
		fmt.Printf("Случайное число: %d\n", rand.Intn(100))
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Генерация случайных чисел завершена!")
}

func sum(arr []int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Сумма начинает вычисляться...")
	time.Sleep(time.Second)
	sum := 0
	for _, v := range arr {
		sum += v
	}
	
	fmt.Printf("Сумма чисел (1, 2, 3, 4, 5) = %d\n", sum)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go genereateRandomNumbers(5, &wg)
	go factorial(4, &wg)
	go sum([]int{1, 2, 3, 4, 5}, &wg)

	wg.Wait()
	fmt.Println("Конец")
}