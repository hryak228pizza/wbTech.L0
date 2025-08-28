package cache

import (
    "database/sql"
    _ "github.com/lib/pq"

	"github.com/hashicorp/golang-lru/v2"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
)

type Cache struct {
	lru *lru.Cache[string, *model.Order] // {uid: Order}
}

// NewCache creates new cache size of N
func NewCache(size int, db *sql.DB) (*Cache, error) {

	// create empty map
	cache, err := lru.New[string, *model.Order](size)
	if err != nil {
		return nil, err
	}
	c := &Cache{ lru: cache }

	// get last N orders 
	lastOrders, err := db.Query("SELECT * FROM orders ORDER BY date_created DESC LIMIT $1", size)
	if err != nil {
		return nil, err
	}
	defer lastOrders.Close()

	for lastOrders.Next() {
        o := &model.Order{}

		// parse order info
        if err := lastOrders.Scan(
            &o.OrderUID, &o.TrackNumber, &o.Entry,
            &o.Locale, &o.InternalSignature, &o.CustomerID,
            &o.DeliveryService, &o.ShardKey, &o.SmID,
            &o.DateCreated, &o.OofShard,
        ); err != nil {
            return nil, err
        }

        // parse delivery
        if err := db.QueryRow(`
            SELECT name, phone, zip, city, address, region, email
            FROM delivery WHERE order_uid = $1
        `, o.OrderUID).Scan(
            &o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip,
            &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
        ); err != nil {
            return nil, err
        }

        // parse payment
        if err := db.QueryRow(`
            SELECT transaction, request_id, currency, provider, amount,
                   payment_dt, bank, delivery_cost, goods_total, custom_fee
            FROM payment WHERE transaction = $1
        `, o.OrderUID).Scan(
            &o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency,
            &o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDT,
            &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
        ); err != nil {
            return nil, err
        }

        // parse items
        itemRows, err := db.Query(`
            SELECT id, chrt_id, track_number, price, rid, name, sale, size,
                   total_price, nm_id, brand, status
            FROM items WHERE order_uid = $1
        `, o.OrderUID)
        if err != nil {
            return nil, err
        }

        for itemRows.Next() {
            it := &model.Item{OrderUID: o.OrderUID}
            if err := itemRows.Scan(
                &it.ID, &it.ChrtID, &it.TrackNumber, &it.Price,
                &it.Rid, &it.Name, &it.Sale, &it.Size,
                &it.TotalPrice, &it.NmID, &it.Brand, &it.Status,
            ); err != nil {
                return nil, err
            }
            o.Items = append(o.Items, it)
        }
        itemRows.Close()

		// save order to cache
        c.SetOrder(o)
    }

	return c, nil
}

// GetOrder returns order by ID and true/false if found/not found
func (c *Cache) GetOrder(id string) (*model.Order, bool) {
	return c.lru.Get(id)	
}

// SetOrder saves order in cache
func (c *Cache) SetOrder(order *model.Order) {
	c.lru.Add(order.OrderUID, order)
}