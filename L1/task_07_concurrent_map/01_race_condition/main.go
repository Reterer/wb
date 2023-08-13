/*
Реализовать конкурентную запись данных в map.
*/
package main

import (
	"l1/task_07_concurrent_map/utils"
	"time"
)

/*
	По ссылке https://go.dev/doc/faq#atomic_maps дается ответ, почему
	map не безопасный для конкурентности.

	Поэтому я сделаю условную структуру, которая придаст небольшой смысл.
	Это будет банк со счетами клиентов.

	Мы можем изменять сумму (через функцию add) и проверять
	Я опустил разные проверки, в данном задании делается акцент
	на работу с конкурентностью
*/

type Bank struct {
	accounts map[string]int // условно Uid -> monies (Я хотел сделать для каждого свой тип, но это перегружет пример)
}

func MakeBank() *Bank {
	return &Bank{
		accounts: make(map[string]int),
	}
}

func (b *Bank) Add(uid string, amount int) {
	b.accounts[uid] += amount
}

func (b *Bank) Check(uid string) int {
	return b.accounts[uid]
}

func main() {
	// Создаем банк
	bank := MakeBank()
	// Здесь мы только считываем значения map в 2 горутины
	// Если map не меняется, то мы можем спокойно
	// считывать ее значения конкуретно
	utils.StressBank(bank, utils.StressBankArgs{
		Timeout:      1 * time.Second,
		AddWorkers:   0,
		CheckWorkers: 2,
		Users:        2,
	})
	// Но записывать или считывать и записывать конкуретно
	// мы не можем.
	// Данная функция вызовет панику из-за попытки изменить поле, которое
	// либо считывается, либо записывается в данный момент
	// utils.StressBank(bank, utils.StressBankArgs{
	// 	Timeout:      1 * time.Second,
	// 	AddWorkers:   2,
	// 	CheckWorkers: 2,
	// 	Users:        2,
	// })
}
