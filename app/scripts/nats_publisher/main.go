package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Reterer/wb/app/config"
	"github.com/Reterer/wb/app/models"

	"github.com/nats-io/stan.go"
)

func setRealyUnique(o models.Order) models.Order {
	r := strconv.Itoa(int(time.Now().UnixMicro()))
	o.Uid += r
	o.Delivery.Name += r
	o.Payment.Transaction = o.Uid

	itemsCount := rand.Intn(5)
	item := o.Items[0]
	newItems := make([]models.Item, itemsCount)
	for i := 0; i < itemsCount; i++ {
		newItems[i] = item
		newItems[i].ChrtId = itemsCount
		newItems[i].TrackNumber = r
	}
	o.Items = newItems

	return o
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

var testmsg = []byte(`{
    "order_uid": "test",
    "track_number": "WBILMTESTTRACK",
    "entry": "WBIL",
    "delivery": {
        "name": "Test Testov",
        "phone": "+9720000000",
        "zip": "2639809",
        "city": "Kiryat Mozkin",
        "address": "Ploshad Mira 15",
        "region": "Kraiot",
        "email": "test@gmail.com"
    },
    "payment": {
        "transaction": "b563feb7b2b84b6test",
        "request_id": "",
        "currency": "USD",
        "provider": "wbpay",
        "amount": 1817,
        "payment_dt": 1637907727,
        "bank": "alpha",
        "delivery_cost": 1500,
        "goods_total": 317,
        "custom_fee": 0
    },
    "items": [
        {
            "chrt_id": 9934930,
            "track_number": "WBILMTESTTRACK",
            "price": 453,
            "rid": "ab4219087a764ae0btest",
            "name": "Mascaras",
            "sale": 30,
            "size": "0",
            "total_price": 317,
            "nm_id": 2389212,
            "brand": "Vivienne Sabo",
            "status": 202
        }
    ],
    "locale": "en",
    "internal_signature": "",
    "customer_id": "test",
    "delivery_service": "meest",
    "shardkey": "9",
    "sm_id": 99,
    "date_created": "2021-11-26T06:22:19Z",
    "oof_shard": "1"
}`)

func breakData(data []byte) []byte {
	return data[:rand.Intn(len(data))]
}

func breakOrder(data []byte) []byte {
	r := rand.Intn(4)
	p := make(map[string]interface{})
	json.Unmarshal(data, &p)
	if r == 0 {
		// Удалим delivery
		delete(p, "delivery")
	} else if r == 1 {
		// Удалим delivery
		delete(p, "payment")
	} else if r == 2 {
		// Удалим items
		delete(p, "items")
	} else if r == 3 {
		// Удалим, например, order_uid
		delete(p, "order_uid")
	}
	data, _ = json.Marshal(p)
	return data
}

func main() {
	// TODO push n random orders
	cfg, err := config.GetConfig()
	cfg.NATS.ClientID = "test-pulisher"
	if err != nil {
		fatalError(err)
	}
	sc, err := stan.Connect(cfg.NATS.ClusterID, cfg.NATS.ClientID)
	if err != nil {
		fatalError(err)
	}

	n := 30  // Количество новых записей
	p := 0.0 // Вероятность сломаной версии (неполные данные) Нет проверок
	q := 0.1 // Вероятность сломать json (мусор)
	var defo models.Order
	_ = json.Unmarshal(testmsg, &defo)

	for i := 0; i < n; i++ {
		var data []byte
		o := setRealyUnique(defo)
		r := rand.Float64()
		data, _ = json.Marshal(o)

		if r < p {
			data = breakOrder(data)
		} else if r < p+q {
			data = breakData(data)
		} else {
		}
		sc.Publish("orders", data)

	}

	defer sc.Close()
}
