package pattern

import "math"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Применимость:
1. Над элементами должны выполняться разнообразные, не связанные между собой методы.
2. Существуют элементы с разными интерфейсами, и мы хотим для них выполнять операции, зависящие от конкретных типов
3. Новые элементы добавляются редко, но новые методы - часто.

Плюсы:
1. Упрощает добавление новых операций.
2. Объединение логики операции для разных типов в одном месте.
3. Есть возможность обходить коллекцию элементов и сохранять внутреннее состояние для посетителя.

Минусы:
1. Сложность добавления новых элементов.
2.  Необходимость реализововать элемент как публичную структуру,
	либо как структуру с развитым интерфейсом. Возможно нарушение инкапсуляции.
*/

/*
Мы знаем, что новые структуры не будут добавляться. Но при этом мы не знаем, какие операции
понадобятся в будущем. Кроме этого, они могут быть соверешенно разнородными.
*/

// Интерфейс посетителя
// Думаю, что благодаря неявной реализации интерфейсов в Go
// мы можем делать много маленьких интерфейсов для Визитеров
// Но здесь я буду придерживаться классической реализации.
type Visitor interface {
	VisitPoint(p *Point)
	VisitCirlce(p *Circle)
	VisitRectangle(p *Rectangle)
}

// Элементы, которые могут принимать посетителя
type Point struct {
	X, Y float64
}

func (p *Point) Accept(v Visitor) {
	v.VisitPoint(p)
}

type Circle struct {
	C Point
	R float64
}

func (c *Circle) Accept(v Visitor) {
	v.VisitCirlce(c)
}

type Rectangle struct {
	A, B Point
}

func (r *Rectangle) Accept(v Visitor) {
	v.VisitRectangle(r)
}

// Теперь сделаем несколько посетителей
// Посетитель для вычисления площади
type SquareVisiter struct {
	sq float64
}

func (v *SquareVisiter) GetSquare() float64 {
	return v.sq
}
func (v *SquareVisiter) VisitPoint(p *Point) {
	v.sq += 0
}
func (v *SquareVisiter) VisitCirlce(p *Circle) {
	v.sq += p.R * p.R * math.Pi
}
func (v *SquareVisiter) VisitRectangle(p *Rectangle) {
	dx := math.Abs(p.B.X - p.A.X)
	dy := math.Abs(p.B.Y - p.A.Y)
	v.sq += dx * dy
}

// Таким образом можно делать много разных посетителей
// Думаю, что в го это даже удобнее, чем в других языках

type Element interface {
	Accept(Visitor)
}

// Клиент, который использует визитера
func CalcSquare(els []Element) float64 {
	sq := &SquareVisiter{}
	for _, el := range els {
		el.Accept(sq)
	}
	return sq.GetSquare()
}
