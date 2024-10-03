package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

func getAreas(s[] Shape) {
	for _, v := range s {
		fmt.Printf("%.2f ", getArea(v))
	}
}

func getArea(s Shape) float64 {
	return s.Area()
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

type Rectangle struct {
	width float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func main() {
	c := Circle{radius: 2.5}
	r := Rectangle{4, 5}
	fmt.Printf("Площадь круга (R = %.2f) = %.2f\n", c.radius, getArea(c))
	fmt.Printf("Площадь прямоугольника (%.2f, %.2f) = %.2f\n\n", r.width, r.height, getArea(r))

	s := []Shape{c, r}
	getAreas(s)
}