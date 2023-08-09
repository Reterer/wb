/*
Реализовать пересечение двух неупорядоченных множеств.
*/
/*
В языке Go нет стандартной структуры множеств
Но можно использовать отображение со значением типа struct{}.
*/

package main

import (
	"fmt"
	"strings"
)

type Set map[string]struct{}

func (s Set) Insert(v string) {
	s[v] = struct{}{}
}

func (s Set) Has(v string) bool {
	_, ok := s[v]
	return ok
}

func (s Set) String() string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func (a Set) Intersection(b Set) Set {
	c := make(Set)

	if len(a) > len(b) {
		// Небольшая оптимизация
		// Мы можем поменять значения, это не изменит внешние переменные
		a, b = b, a
	}

	for k := range a {
		if b.Has(k) {
			c.Insert(k)
		}
	}

	return c
}

func main() {
	a := make(Set)
	a.Insert("A")
	a.Insert("B")
	a.Insert("C")
	a.Insert("D")
	fmt.Println("a: ", a)

	b := make(Set)
	b.Insert("C")
	b.Insert("D")
	b.Insert("E")
	fmt.Println("b: ", b)

	c := a.Intersection(b)
	fmt.Println("c: ", c)
}
