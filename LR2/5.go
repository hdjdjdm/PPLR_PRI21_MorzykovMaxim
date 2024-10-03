package main

import "fmt"

type Rectangle struct {
	width float64
	height float64
}

func (R Rectangle) Square() float64 {
	return R.width * R.height
}

func main() {
	var w, h float64

	fmt.Print("Введите ширину и высоту прямоугольника: ")
	_, err := fmt.Scan(&w, &h)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	
	var rectangle = Rectangle{width: w, height: h}

	fmt.Printf("Площадь прямоугольника: %.2f\n", rectangle.Square())
}