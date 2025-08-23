package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
	model "github.com/hryak228pizza/wbTech.L0/internal/model"
)

const (
	topic         = "new-orders-topic"
	brokerAddress = "localhost:9093"
	groupID       = "new-orders-group"
)

// runs a consumer for reading all incoming orders with kafka
func Consumer(db *sql.DB) {

	// init reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: groupID,
	})
	defer r.Close()

	fmt.Printf("Консьюмер подписан на топик '%s' в группе '%s'\n\n", topic, groupID)

	// TODO:
	// Logging with zap

	ctx := context.Background()
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("ошибка чтения из Kafka: %v\n", err)
		}

		var order model.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("ошибка парсинга JSON: %v\n", err)
			continue
		}

		if err := saveOrder(db, &order); err != nil {
			log.Printf("ошибка записи заказа в БД: %v\n", err)
		} else {
			fmt.Printf("Заказ %s сохранён в БД\n", order.OrderUID)
		}
	}

	// ctx, final := context.WithCancel(context.Background())
	// defer final()

	// readMsg(ctx, r)
}

func saveOrder(db *sql.DB, o *model.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	// TODO:
	// Transactions


	// orders
	_, err = tx.Exec(`INSERT INTO orders 
		(order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		return err
	}

	// delivery
	_, err = tx.Exec(`INSERT INTO delivery 
		(order_uid, name, phone, zip, city, address, region, email) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		return err
	}

	// payment
	_, err = tx.Exec(`INSERT INTO payment 
		(transaction, request_id, currency, provider, amount, payment_dt, 
		bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency,
		o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDT,
		o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal,
		o.Payment.CustomFee)
	if err != nil {
		return err
	}

	// items
	for _, item := range o.Items {
		_, err = tx.Exec(`INSERT INTO items 
			(order_uid, chrt_id, track_number, price, rid, name, sale, size, 
			total_price, nm_id, brand, status) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			o.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}


func readMsg(ctx context.Context, r *kafka.Reader) {
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("ошибка чтения из Kafka: %v\n", err)
		}
		fmt.Printf("Сообщение в топике %v, партиция %v, offset %v: \n\t%s %s\n\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}