/*
Поменять местами два числа без создания временной переменной.
*/
package main

import "fmt"

func main() {

	// 1 (Самый правильный) использовать множественное присваивание
	{
		a, b := 42, 1

		a, b = b, a

		fmt.Println(a, b)
	}

	// А теперь ненормальные
	// 2 Использовать сложение и вычитание
	{
		a, b := 42, 1
		// Возможно переполнение, но алгоритм по идее должен работать. (https://go.dev/ref/spec#Integer_overflow)

		a = a + b // a + b
		b = a - b // a + b - b = a
		a = a - b // a + b - a = b

		fmt.Println(a, b)
	}

	// 3 Использовать xor
	{
		a, b := 42, 1
		a = a ^ b // a ^ b
		b = a ^ b // a ^ b ^ b = a
		a = a ^ b // a ^ b ^ a = b

		fmt.Println(a, b)
	}
}
