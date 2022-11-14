package store

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid" db:"order_uid"`
	Delivery          Delivery  `json:"delivery" db:"-" goqu:"skipinsert"`
	Payment           Payment   `json:"payment" db:"-" goqu:"skipinsert"`
	Items             []Item    `json:"items" db:"-" goqu:"skipinsert"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Locate            string    `json:"locale" db:"locate"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerId        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	ShardKey          string    `json:"shardkey" db:"shardkey"`
	SmId              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

type OrdersGateway interface {
	GetOrder(uid string) (Order, error)
	AddOrder(order Order) error
	GetAllOrders() ([]Order, error)
}

func (db *DB) GetAllOrders() ([]Order, error) {
	orders := make([]Order, 0)
	var err error
	err = db.From("order").Select("order_uid").ScanStructs(&orders)
	if err != nil {
		return orders, err
	}

	for _, order := range orders {
		o, err := db.GetOrder(order.OrderUid)
		if err != nil {
			return orders, err
		}

		orders = append(orders, o)
	}

	return orders, err
}
func (db *DB) GetOrder(uid string) (Order, error) {
	order, b := db.cache[uid]
	if b {
		return order, nil
	}

	order = Order{}
	var err error

	delivery := Delivery{}
	b, err = db.Select().From("delivery").Where(goqu.Ex{"order_id": uid}).ScanStruct(&delivery)
	if err != nil {
		return order, err
	}
	if !b {
		return order, fmt.Errorf("Does not contains")
	}

	items := make([]Item, 0)
	err = db.Select().From("item").Where(goqu.C("order_id").Eq(uid)).ScanStructs(&items)
	if err != nil {
		return order, err
	}

	payment := Payment{}
	b, err = db.Select().From("payment").Where(goqu.C("order_id").Eq(uid)).ScanStruct(&payment)
	if err != nil {
		return order, err
	}
	if !b {
		return order, fmt.Errorf("Does not contains")
	}

	b, err = db.Select().From("order").Where(goqu.C("order_uid").Eq(uid)).ScanStruct(&order)
	if err != nil {
		return order, err
	}
	if !b {
		return order, fmt.Errorf("Does not contains")
	}
	order.Payment = payment
	order.Items = items
	order.Delivery = delivery
	return order, err

}

func (db *DB) AddOrder(order Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	order.Payment.OrderId = order.OrderUid
	order.Delivery.OrderId = order.OrderUid
	sql, _, _ := db.Insert("order").Rows(order).ToSQL()
	fmt.Println(sql)
	_, err = db.Insert("order").Rows(order).Executor().Exec()
	if err != nil {
		return err
	}

	for _, i := range order.Items {
		i.OrderId = order.OrderUid
		_, err = db.Insert("item").Rows(i).Executor().Exec()
		if err != nil {
			return err
		}
	}
	_, err = db.Insert("payment").Rows(order.Payment).Executor().Exec()
	if err != nil {
		return err
	}

	_, err = db.Insert("delivery").Rows(order.Delivery).Executor().Exec()
	if err != nil {
		return err
	}

	err = tx.Commit()

	db.cache[order.OrderUid] = order
	return err
}
