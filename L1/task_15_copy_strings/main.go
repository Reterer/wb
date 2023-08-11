/*
К каким негативным последствиям может привести данный фрагмент кода, и как это исправить?
Приведите корректный пример реализации.
*/

package main

import (
	"fmt"
	"strings"
)

/*
	var justString string
	func someFunc() {
		v := createHugeString(1 << 10)
		justString = v[:100]
	}

	func main() {
		someFunc()
	}

	Реализация сверху с моей точки зрения имеет несколько проблем:
	1.
		Индексация строки происходит в байтах, хотя сама она представлена в
		кодировке utf-8, где символ может иметь размер от 1 до 4 байтов.

		Это означает, что функция может повредить последний символ,
		если в строке используются не только ascii символы.

		И при использовании этой функции - это нужно учитывать.

		В качестве решения можно использовать перевод в срез рун, а затем
		обратно в строку.
		Либо проитерироваться по символам строки.
	2.
		Большая и маленькая строка делят общую память, поэтому GC не может
		освободить помять, которую занимала большая строка, хоть ее большая часть и не используется

		Решить можно с помощью копирования: strings.Clone() (https://pkg.go.dev/strings#Clone)
	3.
		Использование переменной на уровне пакета.
		В будущем это может ухудшить читаемость кода, так как сложно понять, что влияет на эту переменную.

		Стараться не использовать изменяемые переменные на уровне пакета.
*/

func createHugeString(n int) string {
	return strings.Repeat("ф", n)
}

// Копирует первые N БАЙТОВ строки
func copyFirstNBytes(s string, n int) string {
	return strings.Clone(s[:n])
}

// Копирует первые N СИМВОЛОВ (рун) строки
func copyFirstNRunes(s string, n int) string {
	// Вообще, можно не создовать срез рун (это может быть накладно для больших строк)
	// Но так короче и понятнее
	r := []rune(s)
	newStr := string(r[:n])
	return newStr
}

func main() {
	s := createHugeString(1 << 10)
	fmt.Println("copy first 5 bytes:", copyFirstNBytes(s, 5)) // фф�
	fmt.Println("copy first 5 runes:", copyFirstNRunes(s, 5)) // ффффф
}