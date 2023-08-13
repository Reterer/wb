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
	Используем канал как мютекс
	Идея в том, что мы можем сделать канал с размером буффера равным 1.

	Это решение выглядит лучше, чем прошлое.
	Оно проще воспринимается, а еще потребляет меньше ресурсов.
*/

type Bank struct {
	accounts map[string]int
	ch       chan struct{}
}

func MakeBank() *Bank {
	return &Bank{
		accounts: make(map[string]int),
		ch:       make(chan struct{}, 1),
	}
}

func (b *Bank) Add(uid string, amount int) {
	b.ch <- struct{}{} // Блокируем вход в критическую область
	b.accounts[uid] += amount
	<-b.ch // Снимаем блок
}

func (b *Bank) Check(uid string) int {
	b.ch <- struct{}{} // Блокируем вход в критическую область
	amount := b.accounts[uid]
	<-b.ch // Снимаем вход
	return amount
}

func main() {
	// Создаем банк
	bank := MakeBank()

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
