package main

import (
	"fmt"
	"math"
)

type Circle struct {
	radius float64
}

func (c Circle) getArea() float64 {
	return math.Pi * c.radius * c.radius
}

func main() {
	c := Circle{radius: 2.5}
	fmt.Printf("Площадь круга (R = %.2f) = %.2f", c.radius, c.getArea())
}