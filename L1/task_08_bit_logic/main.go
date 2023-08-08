/*
Дана переменная int64. Разработать программу которая устанавливает i-й бит в 1 или 0.
*/

package main

import (
	"fmt"
)

// Возвращает bitset с установленным i-м битом в значение v
// Если v != 0 или v != 1, то тогда функция ничего не делает,
// Поэтому нужно быть осторожным с ней
// В нашем распоряжении имеется 63 бита с нулевого по 62-ой.
// При изменении 63-го бита у нас получается некорректный результат
// Это связанно с тем, что м исползуем знаковый тип
// По-хорошему нужно использовать беззнаковый тип uint64
func setNthBit(bitset int64, i int, v int) int64 {
	bitset &= ^(1 << i)     // Обнуляем i-ый бит
	bitset |= int64(v) << i // Устанавливаем значение нужное значение в i-ый бит
	return bitset
}

func main() {
	var setbit int64
	// установим нулевой, второй и третий бит в 1
	setbit = setNthBit(setbit, 0, 1)
	setbit = setNthBit(setbit, 2, 1)
	setbit = setNthBit(setbit, 3, 1)
	fmt.Printf("%b\n", setbit) // print 1101

	// Установим второй бит в 0
	setbit = setNthBit(setbit, 2, 0)
	fmt.Printf("%b\n", setbit) // print 1001

}