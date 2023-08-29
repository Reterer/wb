package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
	Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
	У нас есть какая-то подсистема/фреймворк, который занимается обработкой выражений
	Это сложная система, где есть много структур.
*/
type Parser struct{}

func (p *Parser) Parse(data []byte) *Node {
	fmt.Println("Parse from Parser")

	return nil
}

type Node struct{}

type Optimizer struct{}

func (o *Optimizer) Optimize(root *Node) *Node {
	fmt.Println("Optimize from Optimizer")

	return root
}

type Executer struct{}

func (e *Executer) Exec(root *Node) *Result {
	fmt.Println("Exec from Executer")

	return &Result{}
}

type Result struct{}

func (r *Result) String() string {
	return "some result"
}

// Фасад, который упрощает взаимодействие с системой, реализующий ограниченные возможности
type CalculatorFacade struct {
	parser    *Parser
	optimizer *Optimizer
	executer  *Executer
}

// Настраиваем копоненты системы и возвращаем новый фасад
func NewCalculatorFacade() *CalculatorFacade {
	return &CalculatorFacade{
		parser:    &Parser{},
		optimizer: &Optimizer{},
		executer:  &Executer{},
	}
}

// За функцией calc скрывается множество операций с системой
func (calc *CalculatorFacade) Calc(expression string) string {
	data := []byte(expression)
	tree := calc.parser.Parse(data)
	tree = calc.optimizer.Optimize(tree)

	res := calc.executer.Exec(tree)
	return res.String()
}
