/*
Разработать программу, которая проверяет,
что все символы в строке уникальные (true — если уникальные, false etc).
Функция проверки должна быть регистронезависимой.

Например:
abcd — true
abCdefAaf — false
aabcd — false
*/
package main

import (
	"fmt"
	"unicode"
)

// Проверяет, что все символы в строке уникальны
// Функция регистро-независимая
func checkUnique(s string) bool {
	// Я хочу, что бы она работала с unicode, поэтому буду использовать map
	dict := make(map[rune]struct{})

	for _, r := range s {
		lowerR := unicode.ToLower(r)
		if _, ok := dict[lowerR]; ok {
			return false // Такой символ уже есть в словаре
		}
		dict[lowerR] = struct{}{}
	}

	return true
}

func main() {
	fmt.Println("abcd", checkUnique("abcd"))
	fmt.Println("abCdefAaf", checkUnique("abCdefAaf"))
	fmt.Println("aabcd", checkUnique("aabcd"))
	fmt.Println("абвг", checkUnique("абвг"))
	fmt.Println("абвгА", checkUnique("абвгА"))
}
