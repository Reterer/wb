package nats

import (
	"encoding/json"
	"fmt"
	"l0/models"

	"github.com/nats-io/stan.go"
)

func MakeOrderHandler(orderModel models.OrderModel) func(*stan.Msg) {
	return func(m *stan.Msg) {
		/*
			1. Распарсить json
			2. Вставить в модель
		*/

		var order models.Order
		// TODO валидация json согласно схеме
		if err := json.Unmarshal(m.Data, &order); err != nil {
			fmt.Printf("info: nats order handler can't unmarshal: %v\n", err)
			return
		}
		// TODO возможна какая-то дополнительная валидация
		if err := orderModel.Insert(order); err != nil {
			fmt.Printf("info: can't insert order: %v | %v\n", order.Uid, err)
			return
		}
	}
}
