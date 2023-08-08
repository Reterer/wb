/*
Реализовать конкурентную запись данных в map.
*/
package main

import (
	"fmt"
	"l1/task_07_concurrent_map/utils"
	"time"
)

/*
	Используем каналы и горутину обработчик,
	Идея заключается в следующем:
		Создаем горутину, которая будет обрабатывать
		команды формата func(map[string]int).
		И вся работа с map должна вестись только через эту функцию.

	Выглядит это, конечно, интересно. Но такую структуру сложно изменять,
	а еще здесь крутится горутина.
*/

type Bank struct {
	cmd chan func(map[string]int) // Очередь команд, которые будут выполняться горутиной
}

func MakeBank() (*Bank, func()) {
	bank := &Bank{
		cmd: make(chan func(map[string]int)), // Создаем очередь команд
	}
	// Создаем горутину, которая будет выполнять команды
	go func() {
		defer fmt.Println("map manager is closed")
		// Создаем отображение в этой функции
		// Это гарантирует то, что с ней будет работать только эта функция
		accounts := make(map[string]int)
		for {
			select {
			case f, ok := <-bank.cmd:
				if !ok {
					return
				}
				// Выполняем команды
				f(accounts)
			}
		}
	}()
	// функцию, которая завершит работу горутины
	cancel := func() {
		close(bank.cmd)
	}

	return bank, cancel
}

func (b *Bank) Add(uid string, amount int) {
	// Функция, которая изменет сумму указанному человеку
	b.cmd <- func(m map[string]int) {
		m[uid] += amount
	}
}

func (b *Bank) Check(uid string) int {
	res := make(chan int) // канал для возрата значения
	b.cmd <- func(m map[string]int) {
		// сохраняем значение
		res <- m[uid]
	}
	return <-res
}

func main() {
	// Создаем банк
	bank, cancel := MakeBank()
	defer cancel() // не забываем закрыть горутину

	// Проверим, что эта система работает
	bank.Add("test", 10)
	fmt.Printf("10 == %d\n", bank.Check("test"))

	// Тест на одновременное чтение и запись
	utils.StressBank(bank, utils.StressBankArgs{
		Timeout:      1 * time.Second,
		AddWorkers:   10,
		CheckWorkers: 10,
		Users:        3,
	})
}
