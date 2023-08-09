/*
Дана последовательность температурных колебаний: -25.4, -27.0 13.0, 19.0, 15.5, 24.5, -21.0, 32.5.
Объединить данные значения в группы с шагом в 10 градусов.
Последовательность в подмножноствах не важна.

Пример: -20:{-25.0, -27.0, -21.0}, 10:{13.0, 19.0, 15.5}, 20: {24.5}, etc.
*/
/*
	Как выбрать ключ?
	1.
		Будем считать, что каждая группа соединяется с целым ключом, кратным keyBase
		и вмещает в себя промежуток:
			(key - keyBase, key], если key <= 0
			[key, key + keyBase), если key > 0


		-20: (-30, -20]
		-10: (-20, -10]
		// Пропущен интервал (-10, 0)
		0:   [0, 10)
		10:  [10, 20)
	2.
		Будем считать, что:
			[key, key + keyBase) для любых key

		-30: (-30, -20]
		-20: (-20, -10]
		-10: (-10, 0)	// из-за кусочности функции
		0:   [0, 10)
		10:  [10, 20)
		Проблем никаких не возникает
*/

package main

import "fmt"

func key(x float32, base int) int {
	if x >= 0 {
		return int(x) / base * base
	}
	// Если число меньше нуля, то смещяем его на base, что бы не было коллизии при x близком к 0
	return (int(x) - base) / base * base
}

func groupTemp(m map[int][]float32, keyBase int, seq []float32) {
	// По очереди записывать температуру в нужную группу. Если подходящей группы нет - создать
	for _, t := range seq {
		k := key(t, keyBase)

		arr, ok := m[k]
		if !ok {
			arr = make([]float32, 0)
		}

		arr = append(arr, t)
		m[k] = arr
	}
}

func checkIntervals(m map[int][]float32, base int) bool {
	fbase := float32(base)
	for _, v := range m {
		min := v[0]
		max := v[0]
		for _, el := range v {
			if min > el {
				min = el
			}
			if max < el {
				max = el
			}
		}
		if max-min > fbase {
			return false
		}
	}
	return true
}

func main() {
	seq := []float32{-10.0, -25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, -1.0, 32.5, 9.0, -7.0, 0.0} // Добавил пограничный случай
	keyBase := 10                                                                                    // Интервал, по которому мы будем группировать элементы
	fmt.Println("seq:", seq)

	m := make(map[int][]float32)
	groupTemp(m, keyBase, seq)
	fmt.Println("groupTemp:", m)
	fmt.Println("intervals is correct:", checkIntervals(m, keyBase))

}
