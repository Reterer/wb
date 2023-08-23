package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Применимость:
1. Нужно конструировать разные продукты
2. Алгоритм создания объекта не должен зависеть от продукта

Плюсы:
1. Позволяет легко создавать разные продукты с помощью одного и того же директора
2. Изолирует код, реализующий конструирование и логику продукта
3. Дает больший контроль над процессом конструирования (Пошаговая сборка продукта)

Минусы:
1. Усложняет код из-за дополнительных интерфейсов и структур
*/

/*
Пусть у нас есть несколько структур, которые отвечают за разное представление отчетов.
Строитель позволит с помощью одного и того же кода строить отчеты в разных форматах.
*/

// Интерфейс Строителя. Определяет функции, которые будут строить продукт
type Builder interface {
	AddText(s string)                  // Добовляет какой-то текст
	AddLink(label string, link string) // Добавляет ссылку на что-то
	AddHeader(h string)                // Добавляет заголовок
}

// Строитель для HTMLProduct
type HTMLBuilder struct {
	p *HTMLProduct
}

func NewHTMLBuilder() *HTMLBuilder {
	return &HTMLBuilder{
		p: &HTMLProduct{},
	}
}
func (b *HTMLBuilder) AddText(s string) {
	par := []byte("<p>" + s + "</p>\n")
	b.p.Buf = append(b.p.Buf, par...)
}
func (b *HTMLBuilder) AddLink(label string, link string) {
	a := []byte(fmt.Sprintf("<a href=\"%s\">%s</a>\n", link, label))
	b.p.Buf = append(b.p.Buf, a...)
}
func (b *HTMLBuilder) AddHeader(s string) {
	h := []byte("<h>" + s + "</h>\n")
	b.p.Buf = append(b.p.Buf, h...)
}
func (b *HTMLBuilder) GetHtml() *HTMLProduct {
	return b.p
}

type HTMLProduct struct{ Buf []byte } // Продукт, который представляет страницу в виде HTML

// Строитель для MarkDown
type MDBuilder struct {
	p *MDProduct
}

func NewMDBuilder() *MDBuilder {
	return &MDBuilder{
		p: &MDProduct{},
	}
}
func (m *MDBuilder) AddText(s string) {
	par := []byte(s + "\n")
	m.p.Buf = append(m.p.Buf, par...)
}
func (m *MDBuilder) AddLink(label string, link string) {
	a := []byte(fmt.Sprintf("[%s](%s)\n", label, link))
	m.p.Buf = append(m.p.Buf, a...)
}
func (m *MDBuilder) AddHeader(s string) {
	h := []byte("# " + s + "\n")
	m.p.Buf = append(m.p.Buf, h...)
}
func (m *MDBuilder) GetMD() *MDProduct {
	return m.p
}

type MDProduct struct{ Buf []byte } // Продукт, который представляет страницу в виде MarkDown
// На самом деле продукт может иметь другой формат. Здесь достаточно простой пример.

// Директор, то есть код, который как-то собирает продукт, используя билдера
// Например здесь может быть генерация какого-то отчета.
func BuildReport(b Builder) {
	b.AddHeader("Header")
	b.AddText("just text.")
	b.AddLink("link", "#")
}

// Клиент, которому нужно построить html отображение
func HTMLClient() {
	htmlBuilder := NewHTMLBuilder()
	BuildReport(htmlBuilder)
	htmlBuilder.GetHtml()
}

// Клиент, которому нужно построить md отображение
func MDClient() {
	mdBuilder := NewMDBuilder()
	BuildReport(mdBuilder)
	mdBuilder.GetMD()
}
