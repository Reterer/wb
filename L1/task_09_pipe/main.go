/*
Разработать конвейер чисел.
Даны два канала: в первый пишутся числа (x) из массива, во второй — результат операции x*2,
после чего данные из второго канала должны выводиться в stdout.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
)

// Создает горутину, которая умножает поступающие числа на 2 и записывает их в выходной канал
func doubleWorker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer fmt.Println("doubleWorker is done")
		defer close(out) // закрываем канал, когда закончим запись
		// горутина будет жить, пока не закроется канал in
		for x := range in {
			out <- x * 2
		}
	}()
	return out
}

func printer(in <-chan int) {
	go func() {
		defer fmt.Println("printer is done")
		for x := range in {
			fmt.Println(x)
		}
	}()
}

func main() {
	arr := []int{0, 1, 2, 4, 6, 8}

	ch := make(chan int) // канал для записи чисел
	// создаем конвейер чисел
	out := doubleWorker(ch)
	printer(out)

	// Пишем числа в канал
	for x := range arr {
		ch <- x
	}
	close(ch) // Закрываем канал, что бы другие горутины завершили работу

	// Ожидаем прерывания
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
