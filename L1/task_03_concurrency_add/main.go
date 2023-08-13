/*
Дана последовательность чисел: 2,4,6,8,10.
Найти сумму их квадратов(2^2+3^2+4^2+...) с использованием конкурентных вычислений.
*/
package main

import "fmt"

func main() {
	arr := []int{2, 4, 6, 8, 10}
	// Создаем канал, по которому горутины будут возвращать значения
	retChan := make(chan int)

	worker := func(i int) {
		retChan <- arr[i] * arr[i] // Возвращаем результат работы
	}

	// Запускаем горутины
	for i := 0; i < len(arr); i++ {
		go worker(i)
	}

	// аккумулируем результат в эту переменную
	sum := 0
	for i := 0; i < len(arr); i++ {
		sum += <-retChan
	}

	fmt.Printf("sum of arr: %d\n", sum) // print 220
}
