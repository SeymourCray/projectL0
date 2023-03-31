package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"projectL0/internal/entity"
)

const (
	getAllQuery = `SELECT o.order_uid, o.track_number, o.entry, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email, ` +
		`p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee, ` +
		`o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard ` +
		`FROM orders o ` +
		`INNER JOIN payments p ON p.order_uid = o.order_uid ` +
		`INNER JOIN deliveries d ON o.order_uid = d.order_uid `
	getItemsQuery = `SELECT i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status ` +
		`FROM items_orders io ` +
		`INNER JOIN items i ON io.chrt_id=i.chrt_id ` +
		`WHERE io.order_uid=$1 `
	insertOrderQuery     = `INSERT INTO orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	insertItemQuery      = `INSERT INTO items VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	insertDeliveryQuery  = `INSERT INTO deliveries (name, phone, zip, city, address, region, email, order_uid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	insertPaymentQuery   = `INSERT INTO payments VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	insertItemOrderQuery = `INSERT INTO items_orders (order_uid, chrt_id) VALUES ($1, $2);`
)

type OrderRepo interface {
	GetOrderByID(OrderUID string) (entity.Order, error)
	InsertOrder(order entity.Order)
}

type OrderRepository struct {
	Connection   *sqlx.DB
	cashedOrders []entity.Order
}

func NewOrderRepository(connection *sqlx.DB) *OrderRepository {
	var orders []entity.Order

	err := connection.Select(&orders, getAllQuery)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, order := range orders {
		err = connection.Select(&order.Items, getItemsQuery, order.OrderUID)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	return &OrderRepository{Connection: connection, cashedOrders: orders}
}

func (p *OrderRepository) InsertOrder(order entity.Order) {
	tx, _ := p.Connection.Begin()

	tx.Exec(insertOrderQuery, order.OrderUID, order.TrackNumber, order.Entry,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	tx.Exec(insertDeliveryQuery, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.OrderUID)

	tx.Exec(insertPaymentQuery, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUID)

	for _, item := range order.Items {
		tx.Exec(insertItemQuery, item.ChartID, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)

		tx.Exec(insertItemOrderQuery, order.OrderUID, item.ChartID)
	}

	err := tx.Commit()
	if err != nil {
		log.Println(err.Error())
		return
	}

	p.cashedOrders = append(p.cashedOrders, order)
	log.Println("new order inserted")
}

func (p *OrderRepository) GetOrderByID(orderUID string) (entity.Order, error) {
	var foundedOrder entity.Order

	err := errors.New("order not found")

	for _, order := range p.cashedOrders {
		if order.OrderUID == orderUID {
			foundedOrder = order
			err = nil
		}
	}

	return foundedOrder, err
}
