package ciclebuf

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	N := 5
	buf := NewBuf(N)
	if buf == nil {
		t.Fatalf("buf can't be nil")
	}
}

func checkPop(t *testing.T, buf *Buf, arr []byte) {
	t.Helper()
	res, ok := buf.Pop()
	if !ok {
		t.Fatal("ok should be true")
	}
	if !reflect.DeepEqual(res, arr) {
		t.Errorf("\ngot:  %v\nwant: %v", res, arr)
	}
}

func TestBufZero(t *testing.T) {
	buf := NewBuf(0)

	A := []byte("a string")

	buf.Push(A)
	buf.Push(A)

	_, ok := buf.Pop()
	if ok == true {
		t.Error("ok should be false")
	}
}
func TestBufOne(t *testing.T) {
	N := 1
	buf := NewBuf(N)

	A := []byte("a string")
	B := []byte("b string")

	buf.Push(A)
	buf.Push(A)
	buf.Push(B)
	checkPop(t, buf, B)
}

func TestBufOverlap(t *testing.T) {
	N := 3
	buf := NewBuf(N)

	A := []byte("a string")
	B := []byte("b string")

	buf.Push(A) // A
	buf.Push(A) // A A
	buf.Push(B) // A A B
	buf.Push(A) // A B A
	buf.Push(B) // B A B

	checkPop(t, buf, B)
	checkPop(t, buf, A)
	checkPop(t, buf, B)

	if buf.len != 0 {
		t.Errorf("buf len != 0: len = %d", buf.len)
	}
}
