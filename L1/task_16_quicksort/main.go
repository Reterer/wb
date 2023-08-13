/*
Реализовать быструю сортировку массива (quicksort) встроенными методами языка.
*/
package main

import (
	"fmt"
	"math/rand"
)

func quicksort(a []int) []int {
	// nil или один элемент
	if len(a) < 2 {
		return a
	}

	// Находим опорный элемент
	left, right := 0, len(a)-1
	pivot := rand.Intn(len(a))

	// Перебрасываем опорный элемент в конец, что бы было проще разделять другие элементы
	a[right], a[pivot] = a[pivot], a[right]

	// Перекидываем элементы меньшие чем опорный в левую часть
	for i := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	// На самом деле, правильное место для опорного элемента - между левой и правой частью
	// Это его отсорированная позиция. Поэтому ставим его на правильное место
	a[left], a[right] = a[right], a[left]

	quicksort(a[:left])   // отдельно сортируем только левую часть
	quicksort(a[left+1:]) // только правую часть

	return a
}

func main() {
	fmt.Println("nil:", quicksort(nil))
	fmt.Println("len = 1:", quicksort([]int{5}))
	fmt.Println("len > 1:", quicksort([]int{5, 3, 7, 6, 5, 4, 4, 6, 2, 8, 7}))
	fmt.Println("reversed:", quicksort([]int{6, 5, 4, 3, 2, 1}))
	fmt.Println("sorted:", quicksort([]int{1, 2, 3, 4, 5, 6}))
}
