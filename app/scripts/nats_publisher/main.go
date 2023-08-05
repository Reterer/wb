package main

import (
	"encoding/json"
	"fmt"
	"l0/config"
	"l0/models"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

func setRealyUniqueUid(o models.Order) models.Order {
	o.Uid += time.Now().Format(time.DateTime)
	return o
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

var testmsg = []byte(`{
    "order_uid": "b563feb7b2b84b6test",
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

func main() {
	// TODO push n random orders
	cfg, err := config.GetConfig()
	cfg.NATS.ClientID = "test-pulisher"
	if err != nil {
		fatalError(err)
	}
	fmt.Println(cfg)

	var o models.Order
	_ = json.Unmarshal(testmsg, &o)
	o = setRealyUniqueUid(o)
	data, _ := json.Marshal(o)

	sc, _ := stan.Connect(cfg.NATS.ClusterID, cfg.NATS.ClientID)
	sc.Publish("orders", data)

	defer sc.Close()
}
