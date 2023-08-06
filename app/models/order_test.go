package models

import (
	"encoding/json"
	"l0/config"
	"strconv"
	"testing"
	"time"

	"github.com/go-test/deep"
)

var OrderJson1 = []byte(`{
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

var dbconf = config.DBConfig{
	DBname: "servicedb",
	User:   "serviceuser",
	Pass:   "servicepassword",
	Host:   "127.0.0.1",
	Port:   "5432",
}

var Order1 = Order{
	Uid:         "b563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "b563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Items: []Item{
		{
			ChrtId:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmId:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}

var OrderManyItems = Order{
	Uid:         "manyItems",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "manyItems",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Items: []Item{
		{
			ChrtId:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
		{
			ChrtId:      1234,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "test2",
			Name:        "item2",
			Sale:        123,
			Size:        "0",
			TotalPrice:  317,
			NmId:        232452345,
			Brand:       "lol",
			Status:      202,
		},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmId:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}

func TestOrderJsonUnMarshal(t *testing.T) {
	got := Order{}
	want := Order1

	if err := json.Unmarshal(OrderJson1, &got); err != nil {
		t.Errorf("Err should be nil, but: %s", err.Error())
	}

	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func skipDB(t *testing.T) {
	// if os.Getenv("CI") != "" {
	// 	t.Skip("Skipping testing in CI environment")
	// }
	if false {
		t.Skip("Skipping tests with db interaction")
	}
}

func TestMakeOrderModel(t *testing.T) {
	skipDB(t)

	c, err := MakeCachedOrderModel(dbconf)
	if err != nil {
		t.Errorf("Err should be nil, but: %s", err.Error())
	}
	defer c.Close()
}

func setRealyUniqueUid(o Order) Order {
	o.Uid += strconv.Itoa(int(time.Now().Unix()))
	return o
}

func TestInsertGetByUid(t *testing.T) {
	skipDB(t)

	c, err := MakeCachedOrderModel(dbconf)
	if err != nil {
		t.Fatalf("Err should be nil, but: %s", err.Error())
	}
	defer c.Close()

	insertedOrder := setRealyUniqueUid(OrderManyItems)
	if err := c.Insert(insertedOrder); err != nil {
		t.Errorf("Err should be nil, but: %s", err.Error())
	}

	selectedOrder, err := c.GetByUid(insertedOrder.Uid)
	if err != nil {
		t.Errorf("Err should be nil, but: %s", err.Error())
	}

	// изначально у нас не назначен dbid
	// но GetByUid устанавливает
	// поэтому их нужно сбить для проверки на равенство
	selectedOrder.Payment.dbId = 0
	selectedOrder.Delivery.dbId = 0
	for i := 0; i < len(selectedOrder.Items); i++ {
		selectedOrder.Items[i].dbId = 0
	}

	if diff := deep.Equal(selectedOrder, &insertedOrder); diff != nil {
		t.Error(diff)
	}
}
