package main

import (
	"context"
	"testing"
	"time"
)

func makeChan() chan interface{} {
	ch := make(chan interface{})
	return ch
}

func checkChan(t *testing.T, ch <-chan interface{}, timeout time.Duration, wantOk bool) {
	wantValue := interface{}(nil)
	t.Helper()
	select {
	case v, ok := <-ch:
		if v != wantValue {
			t.Errorf("got: %v want: %v\n", v, wantValue)
		}
		if ok != wantOk {
			t.Errorf("got: %v want: %v\n", ok, wantOk)
		}
	case <-time.After(timeout):
		t.Fatal("timeout")
	}
}

func TestOr(t *testing.T) {
	ch := []chan interface{}{
		makeChan(),
		makeChan(),
		makeChan(),
	}

	timeout := 10 * time.Millisecond

	done := or(ch[0], ch[1], ch[2])
	close(ch[0])
	checkChan(t, done, timeout, true)
	close(ch[1])
	checkChan(t, done, timeout, true)
	close(ch[2])
	checkChan(t, done, timeout, true)

	checkChan(t, done, timeout, false) // канал done закрыт
}

func TestCtxor(t *testing.T) {
	ch := []chan interface{}{
		makeChan(),
		makeChan(),
		makeChan(),
	}

	timeout := 10 * time.Millisecond
	ctx, cancel := context.WithCancel(context.TODO())

	done := ctxor(ctx, ch[0], ch[1], ch[2])
	close(ch[0])
	checkChan(t, done, timeout, true)
	close(ch[1])
	checkChan(t, done, timeout, true)

	cancel()                           // Отменяем контекст
	checkChan(t, done, timeout, false) // канал done закрыт
}
