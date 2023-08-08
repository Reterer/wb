/*
Написать программу, которая конкурентно рассчитает значение квадратов чисел взятых
из массива (2,4,6,8,10) и выведет их квадраты в stdout.
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	// Испольуется конкретно массив, но это не всегда удобно, потому что
	// 1. Массив не является ссылочным типом
	// 2. Массив имеет фиксированный размер
	// 3. Размер массива входит в его тип
	// Однако, иногда массив использовать правильнее, чем срезы.
	// Например для хранения хэшей, которые имеют один и тот-же размер
	arr := [...]int{2, 4, 6, 8, 10}

	// Делаем анонимную функцию для расчета квадратов и их вывода в stdout
	// Она получает индекс, с которым нужно работать через аргумент, что бы скопировать индекс
	// А обращается к массиву через замыкание
	worker := func(i int) {
		sq := arr[i] * arr[i]
		fmt.Printf("i: %d | sq: %d\n", i, sq)
	}

	for i := 0; i < len(arr); i++ {
		// Запустим для каждого элемента свой воркер конкурентно
		go worker(i)
	}

	// Так как горутины завершаются, когда завершается главная программа
	// Нам нужно поставить задержку, что бы горутины успели выполниться
	// (Лучшим решением было бы использовать wait group и ждать завершения других горутин, но я это сделаю в другой раз)
	time.Sleep(1 * time.Second)

	/* print
	i: 4 | sq: 100
	i: 0 | sq: 4
	i: 3 | sq: 64
	i: 1 | sq: 16
	i: 2 | sq: 36
	*/
}