/*
Реализовать бинарный поиск встроенными методами языка.
*/

package main

import (
	"fmt"
)

// Бинарный поиск числа target в отсортированном по возрастанию срезе a
// Вернет индекс первого элемента среза a, такой что выражение element < target ложно
// Если такого элемента нет, то вернет len(a)
// Аналог std::lower_bound из c++
func binSearch(a []int, target int) int {
	l, r := 0, len(a) // [0, n)

	for l < r {
		mid := l + (r-l)/2
		if a[mid] < target {
			l = mid + 1
		} else {
			r = mid
		}
	}

	return l
}

func main() {
	fmt.Println("search nil", binSearch(nil, 1))
	fmt.Println("search empty", binSearch([]int{}, 1))

	fmt.Println("search [1], 1  | 0 ==", binSearch([]int{1}, 1))
	fmt.Println("search [1], 42 | 1 ==", binSearch([]int{1}, 42))
	fmt.Println("search [1], 0  | 0 ==", binSearch([]int{1}, 0))

	fmt.Println("search [1,4,4,4,4,4], 4 | 1 ==", binSearch([]int{1, 4, 4, 4, 4, 4}, 4))
	fmt.Println("search [1,3,4,4,5,7], 5 | 4 ==", binSearch([]int{1, 3, 4, 4, 5, 7}, 5))

	fmt.Println("search [1,3,4,4,5,7], 7  | 5 ==", binSearch([]int{1, 3, 4, 4, 5, 7}, 7))
	fmt.Println("search [1,3,4,4,5,7], 10 | 6 ==", binSearch([]int{1, 3, 4, 4, 5, 7}, 10))
}
