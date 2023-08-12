/*
Разработать программу нахождения расстояния между двумя точками,
которые представлены в виде структуры Point с инкапсулированными параметрами x,y и конструктором.
*/
package main

import (
	"fmt"
	"math"
)

// Отдельный пакет для работы с Point

type Point struct {
	x, y float64
}

func MakePoint(x, y float64) Point {
	return Point{
		x: x,
		y: y,
	}
}

func Distance(a, b Point) float64 {
	dx := b.x - a.x
	dy := b.y - a.y
	return math.Sqrt(dx*dx + dy*dy)
}

// ---------------------------------

func main() {
	a := MakePoint(0, 0)
	b := MakePoint(1, 1)
	fmt.Println(Distance(a, b)) // print 1.4142135623730951
}
