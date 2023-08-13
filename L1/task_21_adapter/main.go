/*
Реализовать паттерн «адаптер» на любом примере.
*/
package main

import (
	"fmt"
	"strconv"
)

// Клиент как-то хранит результаты измерений
type MeasureCollection struct{}

func (mc *MeasureCollection) AddResult(l Resulter) {
	fmt.Println("add result:", l.Result())
}

// Интерфейс, который используется в AddResult
type Resulter interface {
	Result() string
}

// Допустим, что у нас есть внешняя структура, проводящая эксперимент,
// Результат которого нужно сохранить.
type Experiment struct{}

// Однако она имеет другой формат вывода
func (e *Experiment) Result() int {
	return 42 // Какой-то результат эксперимента
}

// Поэтому сделаем адаптер для Эксперимента, так как мы не можем редактировать внешний код
// Либо Experiment может использоваться где-то еще, и мы не хотим загружать эту структуру
type ExperimentAdapter struct {
	experiment *Experiment
}

// Переводим результаты эксперимента в нужный формат
func (e *ExperimentAdapter) Result() string {
	eRes := e.experiment.Result()
	return strconv.Itoa(eRes)
}

func main() {
	mc := MeasureCollection{}

	e := &Experiment{}                // Эксперимент
	eAdapter := &ExperimentAdapter{e} // Адаптер для эксперимента

	mc.AddResult(eAdapter) // С помощью адаптера можно добавить результат эксперимента в коллекцию
}
