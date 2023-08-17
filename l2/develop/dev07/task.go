package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// Мультиплексор для done каналов
// Возвращает канал done, куда сливаются значения из channels
// Если один из каналов channels закрывается, то канал возвращает nil
// Если все каналы будут закрыты, то канал done тоже закроется
func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup // Если все входящие каналы будут закрыты, то закрываем out

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for v := range ch {
				out <- v // Передаем значение, если пришло какое-то значение
			}
			out <- nil // Передаем nil, если канал закрылся
		}(ch)
	}

	// Горутина, отвечающая за закрытие канала out
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Мультиплексор для done каналов с контекстом ctx
// Возвращает канал done, куда сливаются значения из channels
// Если контекст будет отменен, то мультеплексор закроется, включая канал done
// Если один из каналов channels закрывается, то канал возвращает nil
// Если все каналы будут закрыты, то канал done тоже закроется
func ctxor(ctx context.Context, channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup // Если все входящие каналы будут закрыты, то закрываем out

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for {
				select {
				case v, ok := <-ch:
					out <- v
					if !ok { // Канал закрылся
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	// Горутина, отвечающая за закрытие канала out
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v\n", time.Since(start))
}
