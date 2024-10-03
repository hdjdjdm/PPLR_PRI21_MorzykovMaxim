package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var arr [5]int
	
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(100)
	}
	fmt.Printf("Начальный массив: %d\n", arr)

	slice := arr[0:5]
	slice = append(slice, 1234)
	slice = append(slice[:2], slice[3:]...)
	fmt.Printf("Срез: %d\n", slice)
}