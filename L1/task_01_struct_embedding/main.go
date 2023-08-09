/*
Дана структура Human (с произвольным набором полей и методов).
Реализовать встраивание методов в структуре Action от родительской структуры Human (аналог наследования).
*/
package main

import "fmt"

// Cтруктура, описывающая человека
type Human struct {
	ID   int // Поле, которое называется так же, как и поле в структуре Action
	Name string
	Age  int
}

// Текстовое представление структуры Human
func (h Human) String() string {
	return fmt.Sprintf("Name: %s; Age: %d; ID: %d", h.Name, h.Age, h.ID)
}

// Приветствие человека
func (h Human) Greeting() {
	fmt.Printf("Hi, I'm %s, I'm %d years old\n", h.Name, h.Age)
}

type Action struct {
	Human        // Встравивание Human в структуру Action
	ID    int    // Поле, которое называется так же, как и поле в структуре Human
	Act   string // Текущее Действие
}

// Текстовое представление структуры Action
func (a Action) String() string {
	// К вложенной структуре мы можем обратиться напрямую, написав ее тип.
	// Так как a.Human реализует интерфейс Stringer,
	// то printf вызовет метод String у структуры Human
	return fmt.Sprintf("Human: (%v); Act: %s; ID: %d", a.Human, a.Act, a.ID)
}

func main() {
	// Создадим переменную типа Action
	a := Action{
		// Инициализация вложенной структуры Human
		Human: Human{
			Name: "Egor",
			Age:  22,
			ID:   42,
		},
		Act: "sleeping",
		ID:  1,
	}

	// Покажем, что методы вложенных структур можно вызвать напрямую
	a.Greeting() // print "Hi, I'm Egor, I'm 22 years old"

	// При этом, структуры Human и Action имеют метод String,
	// В таком случае будет вызываться метод структуры Action
	fmt.Println(a.String()) // print "Human: (...); Act: sleeping; ID: 1"
	// И так тоже
	fmt.Println(a) // print --- // ---

	// Для того, что бы вызвать метод из структуры Human - нужно явно указать это
	fmt.Println(a.Human.String()) // Метод String структуры Human

	// Также работает и с полями.
	// Уникальные поля встравивются
	fmt.Println("Name:", a.Name)

	// По умолчанию исопльзуется поле структуры Action
	fmt.Println("Action ID:", a.ID)      // print Action ID: 1
	fmt.Println("Human ID:", a.Human.ID) // print Action ID: 1
}
