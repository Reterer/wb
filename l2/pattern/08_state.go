package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
Применимость:
1. Поведение объекта зависит от его состояния и должно изменяться во время выполнения.
2. В коде операций встречаются состоящие из многих ветвей условыне операторы, в которых выбор ветви зависит от состояния.

Плюсы:
1. Локализация поведения, зависящего от состояния, и деление его на части.
2. Явно выраженные переходы между состояниями.

Минусы:
1. Усложнение кода.
2. Состояния могут знать о других состояниях и порождать кучу взаимосвязей.
*/

/*
Например есть чат бот, где операция выполняется в несколько сообщений.
*/

// Разные состояния:
type StartState struct{}

func (s *StartState) Handle(ctx *UserContext, msg string) {
	// ...
	ctx.SetState(ctx.InProgress) // Переход в другое состояние
}

type InProgress struct{}

func (s *InProgress) Handle(ctx *UserContext, msg string) {
	// ...
	ctx.SetState(ctx.Start) // Переход в другое состояние
}

// Интерфейс, которое определяет состояние
type State interface {
	Handle(ctx *UserContext, msg string)
}

// Контекст, который делегирует вызовы сосотоянию
type UserContext struct {
	Start      State
	InProgress State
	curr       State
}

func NewUserContext() *UserContext {
	res := &UserContext{
		Start:      &StartState{},
		InProgress: &InProgress{},
	}
	res.curr = res.Start
	return res
}

// Делегируем обработчик
func (c *UserContext) Handle(msg string) {
	c.curr.Handle(c, msg)
}

// Устанавливаем новое состояние
func (c *UserContext) SetState(s State) {
	c.curr = s
}
