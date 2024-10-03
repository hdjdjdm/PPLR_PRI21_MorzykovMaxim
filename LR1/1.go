package main

import (
	"fmt"
	"time"
)

func main() {
	current_date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Текущие дата и время: %s\n", current_date)
	fmt.Scan()
}
