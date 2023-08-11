/*
Разработать программу, которая в рантайме способна определить тип переменной: int, string, bool, channel из переменной типа interface{}.
*/
package main

import (
	"fmt"
	"reflect"
)

// Определим тип с помощью %T
// Ничем не отличается от reflect.TypeOf(v).String()
// На самом деле в исходном коде Sprintf идет вызов функции reflect.TypeOf(arg).Stirng()
func getTypeSPrintf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// Узнаем тип с помощью рефлексии
// Здесь я решил возвращать разновидность типа
// С одной стороны это позволит узнать, что v interface{} является каналом
// Но здесь мы возвращаем название не типа, а его вида (int, string, bool, chan, map, struct, ...)
func getTypeReflect(v interface{}) string {
	return reflect.TypeOf(v).Kind().String()
}

// Проверим тип с помощью утверждения типа
func getTypeAssertion(v interface{}) string {
	if _, ok := v.(int); ok {
		return "int"
	} else if _, ok := v.(string); ok {
		return "string"
	} else if _, ok := v.(bool); ok {
		return "bool"
	} else if _, ok := v.(chan interface{}); ok {
		// Не получится определить, что это канал, так как в данном случае мы хотим привести переменную к конкретному типу
		return "chan"
	} else {
		return "unknown"
	}
}

// Проверим тип с помощью переключателя типа
func getTypeSwitch(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	case chan interface{}:
		// Не получится определить, что это канал, так как в данном случае мы хотим проверить, что это конкретный тип
		return "chan"
	default:
		return "unknown"
	}
}

// Проверва работоспособности функций
func check(f func(interface{}) string) {
	// Объявляем переменные разных типов
	var i int
	fmt.Printf("want %s \t| got %s\n", "int", f(i))

	var s string
	fmt.Printf("want %s \t| got %s\n", "string", f(s))

	var b bool
	fmt.Printf("want %s \t| got %s\n", "bool", f(b))

	var chi chan int
	fmt.Printf("want %s \t| got %s\n", "chan", f(chi))

	var chs chan struct{}
	fmt.Printf("want %s \t| got %s\n", "chan", f(chs))

	fmt.Println("---")
}

func main() {
	check(getTypeSPrintf)
	check(getTypeReflect)
	check(getTypeAssertion)
	check(getTypeSwitch)
}
