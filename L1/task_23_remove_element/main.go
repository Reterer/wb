/*
Удалить i-ый элемент из слайса.
*/

package main

import "fmt"

// Удаляет i-ый элемент и сдвигает правую часть слайса
// O(n)
func removeIndex(arr []int, i int) []int {
	return append(arr[:i], arr[i+1:]...)
}

// Удаляет i-ый элемент, но нарушает изменяет элементов в слайсе
// O(1)
func removeIndexUnordered(arr []int, i int) []int {
	last := len(arr) - 1
	arr[i], arr[last] = arr[last], arr[i]
	return arr[:last]
}

func main() {
	fmt.Println(removeIndex([]int{1, 2, 3, 4, 5}, 4)) // [1 2 3 4]
	fmt.Println(removeIndex([]int{1, 2, 3, 4, 5}, 2)) // [1 2 4 5]

	fmt.Println(removeIndexUnordered([]int{1, 2, 3, 4, 5}, 4)) // [1 2 3 4]
	fmt.Println(removeIndexUnordered([]int{1, 2, 3, 4, 5}, 2)) // [1 2 5 4]
}
