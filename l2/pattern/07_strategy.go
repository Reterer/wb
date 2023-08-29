package pattern

import "strings"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Применимость:
1. Наличие множества родственных структур, отличающихся только поведением.
2. Наличие нескольких разновидностей алгоритма.
3. Инкапсуляция алгоритма.
4. Много вариантов поведения, представленных разветленными условными операторами.

Плюсы:
1. Возможность замены стратегий на лету.
2. Отделение алгоритма и его внутренних структур от клиентского кода.

Минусы:
1. Клиенты должны знать о различных стратегиях.
2. Усложнение кода.
*/

/*
Например нам нужно выбрать стратегию разделения строки, в зависимости от настроек
*/

// Разделяет на слова
type FieldsStrategy struct{}

func (s *FieldsStrategy) Split(line string) []string {
	return strings.Fields(line)
}

// Разделяет по байтовому разделителю
type StringSepStrategy struct{ Sep string }

func (s *StringSepStrategy) Split(line string) []string {
	return strings.Split(line, s.Sep)
}

// Интерфейс стретегии
type Splitter interface {
	Split(line string) []string
}

// Контест общем случае он может использоваться для передачи данных стратегиям.
// Здесь он выполняет какую-то работу со строками
type SplitContext struct {
	splitter Splitter
}

func (ctx *SplitContext) SetStrategy(s Splitter) {
	ctx.splitter = s
}

func (ctx *SplitContext) HandleLine(line string) string {
	// ...
	_ = ctx.splitter.Split(line)
	// ...
	return ""
}
