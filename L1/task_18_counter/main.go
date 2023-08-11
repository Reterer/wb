/*
Реализовать структуру-счетчик, которая будет инкрементироваться в конкурентной среде.
По завершению программа должна выводить итоговое значение счетчика.
*/
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	acc int64
}

// Увеличивает счетчик на 1 и возвращает новое значение
func (c *Counter) Inc() int64 {
	return atomic.AddInt64(&c.acc, 1)
}

// Возвращает значение счетчика
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.acc)
}

func main() {
	var counter Counter
	var wg sync.WaitGroup

	n := 1000000                        // Общее число инкрементов
	workers := 10                       // Число горутин, которые будут считать
	step := (n + workers - 1) / workers // Деление положительных чисел с округлением вверх

	wg.Add(workers)
	for rem := 0; rem < n; rem += step {
		steps := step
		if rem+steps > n {
			steps = n - rem
		}

		go func(steps int) {
			defer wg.Done()
			for i := 0; i < steps; i++ {
				counter.Inc()
			}
		}(steps)
	}

	wg.Wait()
	fmt.Printf("want: %d got: %d", n, counter.Get())
}
