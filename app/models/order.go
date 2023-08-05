package models

import (
	"database/sql"
	"fmt"
	"l0/config"

	_ "github.com/lib/pq"
)

type Order struct {
	dbId              int64    `json:"-"`
	Uid               string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmId              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	dbId    int64  `json:"-"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	dbId         int64  `json:"-"`
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	dbId        int64  `json:"-"`
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type CachedOrderModel struct {
	db    *sql.DB
	cache map[string]Order
}

// Тут, возможно, не очень правильно делаю
// Но это нужно для удобного тестирования других компонентов
type OrderModel interface {
	Insert(Order) error
	GetByUid(string) (*Order, error)
	Close()
}

func MakeCachedOrderModel(cfg config.DBConfig) (OrderModel, error) {
	connStr := fmt.Sprintf("sslmode=disable host=%s port=%s user=%s password=%s dbname=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DBname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	model := CachedOrderModel{
		db:    db,
		cache: make(map[string]Order),
	}

	// TODO INIT CACHE

	return &model, nil
}

func (c *CachedOrderModel) Insert(order Order) error {
	/*
		1. Вставить основную информацию об ордере
		2. Получить id нового ордера
		3. Вставить Оплату и Доставку
		4. Вставить Товары

		5. Сохранить в кэш
	*/

	// TODO обернуть в транзацкию
	// TODO проверить на дубликаты

	// Вставка основной информации об заказе и получение id записи в БД
	qInsertOrder := `INSERT INTO orders(
		order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	err := c.db.QueryRow(qInsertOrder,
		order.Uid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	).Scan(&order.dbId)
	if err != nil {
		return err
	}

	// Вставка Доставки
	qInsertDelivery := `INSERT INTO deliveries(
		order_id, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = c.db.Exec(qInsertDelivery,
		order.dbId,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		return err
	}

	// Вставка Оплаты
	qInsertPayment := `INSERT INTO payments(
		order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = c.db.Exec(qInsertPayment,
		order.dbId,
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		return err
	}

	// Вставка Товаров
	qInsertItem := `INSERT INTO items(
		order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		// TODO Делать меньше запросов вставки
		_, err = c.db.Exec(qInsertItem,
			order.dbId,
			item.ChrtId,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status)
		if err != nil {
			return err
		}
	}

	// TODO вставить в кеш

	return nil
}

func (c *CachedOrderModel) GetByUid(uid string) (*Order, error) {
	/*
		0. Проверить в кэше

		1. Получить информацию о заказе, доставке, оплате
		2. Получить информацию о товарах
		Считать в Order

		5. Сохранить в кэш
	*/
	// TODO Проверка в кэше
	// TODO транзакция?
	var order Order

	// Заказ
	// Я не использовал INNER JOIN, потому что не хочу всё в одну кучу мешать.
	// В будущем можно будет заменить, что бы делать меньше запросов.

	qSelectOrder := `SELECT id, order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1`
	err := c.db.QueryRow(qSelectOrder, uid).Scan(
		&order.dbId,
		&order.Uid,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard)
	if err != nil {
		return nil, err
	}

	// Доставка
	qSelectDelivery := `SELECT id, name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_id = $1`
	err = c.db.QueryRow(qSelectDelivery, order.dbId).Scan(
		&order.Delivery.dbId,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email)
	if err != nil {
		return nil, err
	}

	// Оплата
	qSelectPayment := `SELECT id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments WHERE order_id = $1`
	err = c.db.QueryRow(qSelectPayment, order.dbId).Scan(
		&order.Payment.dbId,
		&order.Payment.Transaction,
		&order.Payment.RequestId,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	// Товары
	qSelectItems := `SELECT id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_id = $1`
	rows, err := c.db.Query(qSelectItems, order.dbId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	order.Items = make([]Item, 0)
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.dbId,
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	// TODO добавить в кэш
	return &order, nil
}

func (c *CachedOrderModel) Close() {
	if c != nil {
		c.db.Close()
	}
}
