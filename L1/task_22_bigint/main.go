/*
Разработать программу, которая перемножает, делит, складывает, вычитает две числовых переменных a,b, значение которых > 2^20.
*/
package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	{
		// Если операнды вмещаются в тип int32, то
		// Их можно привести к типу int64 и производить с ними нужные операции без переполнения
		a, b := math.MaxInt32, math.MaxInt32
		c := int64(a) * int64(b) // Имеет тип int64
		fmt.Println("mul int32", c)
	}

	{
		// Если числа больше максимально возможного значения int64, то
		// нужно использовать длинную арифметику
		a, b := big.NewInt(1<<62), big.NewInt(1<<62)

		fmt.Println("add", big.NewInt(0).Add(a, b))
		fmt.Println("sub", big.NewInt(0).Sub(a, b))
		fmt.Println("mul", big.NewInt(0).Mul(a, b))
		fmt.Println("div", big.NewInt(0).Div(a, b))
	}
}
