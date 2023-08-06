package models

import (
	"database/sql"
	"errors"
	"fmt"
	"l0/config"

	_ "github.com/lib/pq"
)

type Order struct {
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
	cache map[string]*Order
}

type OrderModel interface {
	Insert(Order) error
	GetByUid(string) (*Order, error)
	ListOfUids() []string
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
		cache: make(map[string]*Order),
	}

	// Будем загружать все записи в хеш
	if err := model.restoreCacheFromDB(); err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *CachedOrderModel) restoreCacheFromDB() error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	qSelectOrders := `SELECT o.order_uid, o.track_number, o.entry, o.locale, 
		o.internal_signature, o.customer_id, o.delivery_service, 
		o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		
		d.id, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		
		p.id, p.transaction, p.request_id, p.currency, p.provider, p.amount, 
		p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
	FROM orders as o INNER JOIN 
		deliveries as d ON d.order_uid = o.order_uid INNER JOIN 
		payments as p ON p.order_uid = o.order_uid`

	oRows, err := tx.Query(qSelectOrders)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer oRows.Close()
	for oRows.Next() {
		order := &Order{
			Items: make([]Item, 0),
		}
		err := oRows.Scan(
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
			&order.OofShard,
			&order.Delivery.dbId,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
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
			&order.Payment.CustomFee,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
		c.cache[order.Uid] = order
	}

	qSelectItems := `SELECT order_uid, id, chrt_id, track_number, price, 
		rid, name, sale, size, total_price, nm_id, brand, status
	FROM items`
	iRows, err := tx.Query(qSelectItems)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer iRows.Close()
	for iRows.Next() {
		var item Item
		var orderUid string
		err := iRows.Scan(
			&orderUid,
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
			tx.Rollback()
			return err
		}
		order := c.cache[orderUid]
		order.Items = append(order.Items, item)
	}
	tx.Commit()
	return nil
}

func (c *CachedOrderModel) Insert(order Order) error {
	/*
		1. Вставить основную информацию об ордере
		2. Получить id нового ордера
		3. Вставить Оплату и Доставку
		4. Вставить Товары

		5. Сохранить в кэш
	*/

	if _, ok := c.cache[order.Uid]; ok {
		return errors.New("such an element already exists")
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	// Вставка основной информации об заказе и получение id записи в БД
	qInsertOrder := `INSERT INTO orders(
		order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(qInsertOrder,
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
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка Доставки
	qInsertDelivery := `INSERT INTO deliveries(
		order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(qInsertDelivery,
		order.Uid,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка Оплаты
	qInsertPayment := `INSERT INTO payments(
		order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(qInsertPayment,
		order.Uid,
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
		tx.Rollback()
		return err
	}

	// Вставка Товаров
	qInsertItem := `INSERT INTO items(
		order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		// TODO Делать меньше запросов вставки
		_, err = tx.Exec(qInsertItem,
			order.Uid,
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
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	// TODO возможно здесь лучше копировать
	c.cache[order.Uid] = &order

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

	if order, ok := c.cache[uid]; ok {
		return order, nil
	}

	var order Order
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	// Заказ
	qSelectOrder := `SELECT o.order_uid, o.track_number, o.entry, o.locale, 
		o.internal_signature, o.customer_id, o.delivery_service, 
		o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		
		d.id, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		
		p.id, p.transaction, p.request_id, p.currency, p.provider, p.amount, 
		p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
	FROM orders as o INNER JOIN 
		deliveries as d ON d.order_uid = o.order_uid INNER JOIN 
		payments as p ON p.order_uid = o.order_uid
	WHERE o.order_uid = $1`
	err = tx.QueryRow(qSelectOrder, uid).Scan(
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
		&order.OofShard,
		&order.Delivery.dbId,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
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
		tx.Rollback()
		return nil, err
	}

	// Товары
	qSelectItems := `SELECT id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1`
	rows, err := tx.Query(qSelectItems, order.Uid)
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	tx.Commit()

	c.cache[order.Uid] = &order
	return &order, nil
}

func (c *CachedOrderModel) ListOfUids() []string {
	orderUids := make([]string, 0, len(c.cache))
	for k := range c.cache {
		orderUids = append(orderUids, k)
	}

	return orderUids
}

func (c *CachedOrderModel) Close() {
	if c != nil {
		c.db.Close()
	}
}
