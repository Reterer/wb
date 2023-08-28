package pattern

import (
	"fmt"
	"log"
)

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Применимость:
1. Запрос может быть обработан боолее чем одним объектом. Объекты должны быть найдены автоматически.
2. Запрос должен быть отправлен одному из нескольких объектов, без явного указания, какому именно.
3. Набор объектов, способных обработать запрос, должен задаваться динамически.

Плюсы:
1. Уменьшает зависимость между клиентом и обработчиками.
2. Гибкость при распределении обязанностей.

Минусы:
1. Запрос может остаться никем не обработанным.
2. Нужно вводить дополнительные объекты-запросы либо формат представления запросов, если они требуют аргументов.
*/

/*
Шаблон "цепочка вызовов" часто используется для обработки каких-либо событий.
Например код, который обрабатывает какой-нибудь пользовательский запрос.
*/

// Запрос
type Message struct {
	From string
}

// Интерфейс для обработчика запроса
type Handler interface {
	Handle(m Message)
}

// Обработчик, который логгирует сообщение
type LogHandler struct {
	Next Handler
}

func (h *LogHandler) Handle(m Message) {
	log.Printf("message from: %s\n", m.From)
	h.Next.Handle(m)
}

// Обработчик, который фильтрует запросы
type UserFilterHandler struct {
	Users []string
	Next  Handler
}

func (h *UserFilterHandler) Handle(m Message) {
	for _, u := range h.Users {
		if m.From == u {
			h.Next.Handle(m)
			return
		}
	}
}

// Обработчик, который "обрабатывает запрос"
type HelloHandler struct{}

func (h *HelloHandler) Handle(m Message) {
	fmt.Printf("Hello, %s!\n", m.From)
}

// Делаем цепочку
func MakeChain() Handler {
	return &LogHandler{
		Next: &UserFilterHandler{
			Users: []string{"Egor"},
			Next:  &HelloHandler{},
		},
	}
}

// Клиент
func ChainOfRespClient() {
	handler := MakeChain()

	// При событии - вызываем обработчик
	handler.Handle(Message{"Egor"})
	handler.Handle(Message{"Andrey"})
}
