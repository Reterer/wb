/*
Разработать программу, которая переворачивает подаваемую на ход строку (например: «главрыба — абырвалг»).
Символы могут быть unicode.
*/
package main

import "fmt"

func reverseChars(s string) string {
	r := []rune(s)
	// Меняем местами руны
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func main() {
	fmt.Println(reverseChars(""))
	fmt.Println(reverseChars("ф"))
	fmt.Println(reverseChars("главрыба"))
}
