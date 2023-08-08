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
	Используем RWMutex для блокировки критических областей
	RWMutex имеет отдельную блокировку для чтения и для записи.
	1. RLock и RUnlock
		Исопльзуются для блокировки критических областей,
		Где данные не меняются.
		Таким образом несколько горутин может считывать значения
		Но при этом критические области, огражденные обычным Lock Unlock
		Блокируются
	2. Lock Unlock - Блокирует и области для чтения и области на изменение
		Следует использовать в областях, где данные меняются.
*/

type Bank struct {
	accounts map[string]int
	l        sync.RWMutex
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
	b.l.RLock() // Блокируем вход в критическую область
	amount := b.accounts[uid]
	b.l.RUnlock() // Снимаем блок
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
