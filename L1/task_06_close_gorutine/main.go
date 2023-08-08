/*
Реализовать все возможные способы остановки выполнения горутины.
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func forRange() {
	data := make(chan int) // Канал для передачи данных
	defer close(data)      // При выходе из функции закрываем канал и высвобождаем горутину

	go func() {
		defer fmt.Println("forRange is done")
		// Цикл завершится, когда канал закроется
		for num := range data {
			fmt.Println("forRange data:", num)
		}
	}()

	for i := 0; i < 1; i++ {
		data <- i
	}
}

func vOk() {
	data := make(chan string)
	defer close(data)

	go func() {
		defer fmt.Println("vOk is done")
		for {
			num, ok := <-data
			// Выходим, если канал закрыт
			if !ok {
				return
			}
			fmt.Println("vOk data:", num)
		}
	}()

	data <- "hi"
}

func done() {
	data := make(chan string)
	// defer close(data) // По хорошему этот канал нужно закрыть, но для показательности не делаем этого

	done := func() chan struct{} {
		done := make(chan struct{}) // Специальный канал, который нужно будет закрыть
		go func() {
			defer fmt.Println("done is done")
			for {
				select {
				case s, ok := <-data:
					// Если канал закрыт, то select начинает из него считывать нулевые значения
					// И если для done это хорошо, то для data - плохо.
					if !ok {
						// Поэтому мы можем ему присвоить nil
						// Чтение из nil всегда является блокирующим (так же как и запись)
						data = nil
						// В данном случае лучше сразу завершить работу горутины
						return
						// Но если бы здесь было бы больше case выражений, то можно использовать continue
					}
					fmt.Println("done data:", s)

				// Когда канал закрывается, из него можно считать нулевое значение сколько угодно раз
				// Поэтому при первой же возможности мы перейдем сюда
				case <-done:
					return
				}
			}
		}()
		return done
	}()

	data <- "info"
	close(done)
}

func cancel() {
	data := make(chan string)
	// defer close(data) // По хорошему этот канал нужно закрыть, но для показательности не делаем этого

	cancel := func() func() {
		done := make(chan struct{})
		go func() {
			defer fmt.Println("cancel is done")
			for {
				select {
				case s, ok := <-data:
					if !ok {
						data = nil
						return
					}
					fmt.Println("cancel data:", s)
				case <-done:
					return
				}
			}
		}()
		// Возвращаем функцию, которая закроет канал done, тем самым завершит работу горутины
		return func() {
			close(done)
		}
	}()

	data <- "something"
	cancel() // Вызываем cancel
}

func ctx() {
	data := make(chan string)
	// Создаем контекст с функцией отмены
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		defer fmt.Println("context is done")
		for {
			select {
			case s, ok := <-data:
				if !ok {
					data = nil
					return
				}
				fmt.Println("context data:", s)
			// По сути все тоже самое, но только обернуто в контекст
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	data <- "something"
	cancel() // Вызываем cancel
}

func main() {
	// 1. Закрыть канал, с помощью которого горутина принимает сообщения
	forRange() // for v := range ch
	time.Sleep(10 * time.Millisecond)

	vOk() // v, ok := <-ch
	time.Sleep(10 * time.Millisecond)

	// 2. Закрыть специальный канал, который горутина "слушает" и завершает работу после его закрытия
	done() // close(done)
	time.Sleep(10 * time.Millisecond)

	// 3. Функция cancel
	// Вместо канала done мы можем вернуть функцию cancel, при вызове которой освободяться ресурсы.
	// У нее есть несколько плюсов, в отличии от канала done:
	// 1. Позволяет скрыть канал done (никто не знает, что могут сделать с каналом в внешнем мире, а функцию можно только вызвать)
	// 2. Может иметь дополнительную логику
	cancel() // cancel()
	time.Sleep(10 * time.Millisecond)

	// 4. Контекст
	// Использует паттерн cancel, однако имеет большую функциональность
	// Контексты могут иметь иерархическую структуру
	// При отмене контекста, все производные от него контексты отменяются
	// Кроме этого, контексты могут сохранять параметры и передавать их по ирерахии.
	ctx()
	time.Sleep(10 * time.Millisecond)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
