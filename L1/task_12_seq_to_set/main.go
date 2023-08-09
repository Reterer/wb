/*
Имеется последовательность строк - (cat, cat, dog, cat, tree) создать для нее собственное множество.
*/

package main

import (
	"fmt"
	"strings"
)

type StringSet map[string]struct{}

func (s StringSet) String() string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func (s StringSet) InsertSlice(l []string) {
	for _, el := range l {
		s[el] = struct{}{}
	}
}

func main() {
	lines := []string{"cat", "cat", "dog", "cat", "tree"}

	set := make(StringSet)
	set.InsertSlice(lines)

	fmt.Println("set of strings:", set)
}
