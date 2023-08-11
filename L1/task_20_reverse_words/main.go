/*
Разработать программу, которая переворачивает слова в строке.
Пример: «snow dog sun — sun dog snow».
*/
package main

import (
	"fmt"
	"strings"
)

func reverseWordsJoin(s string) string {
	words := strings.Fields(s) // split по пробельным символам

	// Меняем местами слова
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return strings.Join(words, " ")
}

func main() {
	fmt.Println(reverseWordsJoin("snow dog sun"))
}
