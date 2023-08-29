Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

**Ответ:**

Сначала будут выведены все элементы, но не в предсказуемом порядке:
```
1
2
.
.
.
7
8
```
Затем, каналы a и b будут закрыты, однако канал c - не будет. горутина Merge будет нагружать поток, на котором работает. В Select будет срабатывать одна из веток и считывать нулевое значение из закрытого канала. А горутина main будет считывать их и выводить нули.