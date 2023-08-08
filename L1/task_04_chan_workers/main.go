/*
Реализовать постоянную запись данных в канал (главный поток).
Реализовать набор из N воркеров, которые читают произвольные данные из канала и выводят в stdout.
Необходима возможность выбора количества воркеров при старте.

Программа должна завершаться по нажатию Ctrl+C.
Выбрать и обосновать способ завершения работы всех воркеров.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func runWorkers(count int) (chan int, chan struct{}, *sync.WaitGroup) {
	var wg sync.WaitGroup       // wait group, которая дает обратную связь о заврешении работы воркеров
	done := make(chan struct{}) // Канал, который уведомляет воркеры о необходимости завершить работу
	data := make(chan int)      // Канал из которого будут считать данные

	for i := 0; i < count; i++ {
		wg.Add(1) // Увеличиваем счетчик
		go func(workerNumber int) {
			defer wg.Done() // Уменьшаем счетчик
			for {
				select {
				case n := <-data:
					// Симулируем сложную задачу
					time.Sleep(10 * time.Millisecond)
					fmt.Printf("worker: %d | got data: %v\n", workerNumber, n)
				case <-done:
					fmt.Printf("worker: %d has completed the work\n", workerNumber)
					return
				}
			}
		}(i)
	}

	return data, done, &wg
}

func main() {
	workersCount := 8 // Кол-во воркеров
	maxI := -1        // Количество записей в канал данных (-1 == очень больше количество записей)

	// Создаем воркеры
	data, done, wg := runWorkers(workersCount)

	// Создаем канал для прерываний
	interrupt := make(chan os.Signal)
	// Подписываемся на прерывания
	signal.Notify(interrupt, os.Interrupt)

	i := 0 // Условные данные, которые будем выводить в stdout
loop: // Почти как goto, но немного лучше...
	for {
		select {
		// Отправляем данные, когда один из воркеров освободится
		case data <- i:
			i++
			if i == maxI {
				fmt.Println("maxI limit")
				break loop
			}
		// Если случилось прерывание, то нужно его обработать
		case sig := <-interrupt:
			fmt.Printf("Got signal: %v\n", sig)
			break loop
		}
	}

	// Закрываем канал done, что бы выполнился второй case у воркеров, который их закроет.
	close(done)
	// При этом, нужно дождаться корректного завершения работы воркеров
	wg.Wait() // Блокирует горутину, пока счетчик не станет равным нулю
}
