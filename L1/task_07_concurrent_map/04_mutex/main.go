/*
Реализовать конкурентную запись данных в map.
*/
package main

import (
	"fmt"
	"l1/task_07_concurrent_map/utils"
	"sync"
	"time"
)

/*
	Используем mutex для блокировки критических областей
	Отличие от канала минимальное
	Однако код с ними выглядет понятнее
	Кроме этого, они работают быстрее, чем каналы
	(https://github.com/SUN-XIN/go-channel-mutex-benchmark)
*/

type Bank struct {
	accounts map[string]int
	l        sync.Mutex
}

func MakeBank() *Bank {
	return &Bank{
		accounts: make(map[string]int),
	}
}

func (b *Bank) Add(uid string, amount int) {
	b.l.Lock() // Блокируем вход в критическую область
	b.accounts[uid] += amount
	b.l.Unlock() // Снимаем блок
}

func (b *Bank) Check(uid string) int {
	b.l.Lock() // Блокируем вход в критическую область
	amount := b.accounts[uid]
	b.l.Unlock() // Снимаем блок
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
