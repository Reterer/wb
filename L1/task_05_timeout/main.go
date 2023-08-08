/*
Разработать программу, которая будет последовательно отправлять значения в канал,
а с другой стороны канала — читать. По истечению N секунд программа должна завершаться.
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	timeout := time.After(1 * time.Second) // Вернет канал, в который через промежуток timout придет метка времени
	data := make(chan int)

	// sender
	go func() {
		for {
			data <- 42
			fmt.Println("send")
		}
	}()

	// reciever
	go func() {
		for {
			<-data
			fmt.Println("recv")
		}
	}()

	<-timeout // Ожидаем нужное количество времени
}
