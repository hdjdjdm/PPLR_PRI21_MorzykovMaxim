package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Operation string
	Num1      float64
	Num2      float64
	Result    chan float64
}

func calculator(reqs <-chan Request) {
	for req := range reqs {
		var result float64
		var err error

		switch req.Operation {
		case "+":
			result = req.Num1 + req.Num2
		case "-":
			result = req.Num1 - req.Num2
		case "*":
			result = req.Num1 * req.Num2
		case "/":
			if req.Num2 == 0 {
				err = fmt.Errorf("деление на ноль")
			} else {
				result = req.Num1 / req.Num2
			}
		default:
			err = fmt.Errorf("неизвестная операция: %s", req.Operation)
		}

		if err != nil {
			fmt.Println("Ошибка:", err)
			req.Result <- 0
		} else {
			req.Result <- result
		}
	}
}

func main() {
	reqs := make(chan Request)

	go calculator(reqs)

	tests := []string{
		"10 + 20",
		"30 - 10",
		"5 * 6",
		"40 / 2",
		"50 / 0",
		"15 ^ 3",
	}

	for _, operation := range tests {
		parts := strings.Fields(operation)
		if len(parts) != 3 {
			fmt.Println("Некорректный ввод:", operation)
			continue
		}

		num1, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			fmt.Println("Ошибка преобразования:", err)
			continue
		}

		num2, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			fmt.Println("Ошибка преобразования:", err)
			continue
		}

		resultChan := make(chan float64)
		req := Request{
			Operation: parts[1],
			Num1:      num1,
			Num2:      num2,
			Result:    resultChan,
		}

		reqs <- req

		result := <-resultChan
		fmt.Printf("%s = %.2f\n", operation, result)
	}

	close(reqs)
}