/*
Реализовать собственную функцию sleep.
*/
package main

import (
	"time"
)

func sleep(duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C
}

// Не блокирует горутину, поэтому максимально грузит поток
// func sleepFor(duration time.Duration) {
// 	start := time.Now()
// 	for time.Since(start) < duration {
// 	}
// }

func main() {
	sleep(2 * time.Second)
}
