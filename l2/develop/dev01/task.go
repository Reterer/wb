package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	ntpServer := "0.europe.pool.ntp.org"
	res, err := ntp.Query(ntpServer)
	if err != nil {
		log.Fatal(err) // По-умолчанию пишет в os.Stdout и завершает программу с кодом ошибки 1 (os.Exit(1))
	}
	fmt.Println("Offset", res.ClockOffset)                               // Значение отстования времени на моем компьютере с временем на сервере
	fmt.Println("Time.Now with offset", time.Now().Add(res.ClockOffset)) // Точное время

	time, err := ntp.Time(ntpServer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current time", time)
}
