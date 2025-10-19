package repository

import (
	"context"
	"database/sql"
	"time"

	sqlc "github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/gen"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
)

type orderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *orderRepo {
	return &orderRepo{DB: db}
}

// saves order in DB
func (r *orderRepo) Save(ctx context.Context, o *model.Order) error {

	// begin transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// inserts:
	// orders
	_, err = tx.ExecContext(ctx, `INSERT INTO orders 
		(order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		return err
	}

	// delivery
	_, err = tx.ExecContext(ctx, `INSERT INTO delivery 
		(order_uid, name, phone, zip, city, address, region, email) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		return err
	}

	// payment
	_, err = tx.ExecContext(ctx, `INSERT INTO payment 
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
		_, err = tx.ExecContext(ctx, `INSERT INTO items 
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

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// conv functions
func str(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func i64(ni sql.NullInt64) *int64 {
	if ni.Valid {
		v := ni.Int64
		return &v
	}
	return nil
}

func tm(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// map functions
func mapDelivery(d sqlc.Delivery, uid string) *model.Delivery {
	return &model.Delivery{
		OrderUID: uid,
		Name:     d.Name.String,
		Phone:    d.Phone.String,
		Zip:      d.Zip.String,
		City:     d.City.String,
		Address:  d.Address.String,
		Region:   d.Region.String,
		Email:    d.Email.String,
	}
}

func mapPayment(p sqlc.Payment) *model.Payment {
	return &model.Payment{
		Transaction:  p.Transaction,
		RequestID:    str(p.RequestID),
		Currency:     str(p.Currency),
		Provider:     str(p.Provider),
		Amount:       i64(sql.NullInt64{Int64: int64(p.Amount.Int32), Valid: p.Amount.Valid}),
		PaymentDT:    i64(p.PaymentDt),
		Bank:         str(p.Bank),
		DeliveryCost: i64(sql.NullInt64{Int64: int64(p.DeliveryCost.Int32), Valid: p.DeliveryCost.Valid}),
		GoodsTotal:   i64(sql.NullInt64{Int64: int64(p.GoodsTotal.Int32), Valid: p.GoodsTotal.Valid}),
		CustomFee:    i64(sql.NullInt64{Int64: int64(p.CustomFee.Int32), Valid: p.CustomFee.Valid}),
	}
}

func mapItems(i []sqlc.Item, uid string) []*model.Item {

	items := make([]*model.Item, 0, len(i))

	for _, item := range i {
		items = append(items, &model.Item{
			ID:          int(item.ID),
			OrderUID:    uid,
			ChrtID:      i64(item.ChrtID),
			TrackNumber: item.TrackNumber,
			Price:       i64(sql.NullInt64{Int64: int64(item.Price.Int32), Valid: item.Price.Valid}),
			Rid:         str(item.Rid),
			Name:        str(item.Name),
			Sale:        i64(sql.NullInt64{Int64: int64(item.Sale.Int32), Valid: item.Sale.Valid}),
			Size:        str(item.Size),
			TotalPrice:  i64(sql.NullInt64{Int64: int64(item.TotalPrice.Int32), Valid: item.TotalPrice.Valid}),
			NmID:        i64(item.NmID),
			Brand:       str(item.Brand),
			Status:      i64(sql.NullInt64{Int64: int64(item.Status.Int32), Valid: item.Status.Valid}),
		})
	}

	return items
}

// helper for mapping orders from sqlc.model to model
func MapToOrder(o sqlc.Order, d sqlc.Delivery, p sqlc.Payment, items []sqlc.Item) *model.Order {

	return &model.Order{
		OrderUID:    o.OrderUid,
		TrackNumber: o.TrackNumber,
		Entry:       str(o.Entry),

		Delivery: *mapDelivery(d, o.OrderUid),
		Payment:  *mapPayment(p),
		Items:    mapItems(items, o.OrderUid),

		Locale:            str(o.Locale),
		InternalSignature: str(o.InternalSignature),
		CustomerID:        str(o.CustomerID),
		DeliveryService:   str(o.DeliveryService),
		ShardKey:          str(o.Shardkey),
		SmID:              i64(sql.NullInt64{Int64: int64(o.SmID.Int32), Valid: o.SmID.Valid}),
		DateCreated:       tm(o.DateCreated),
		OofShard:          str(o.OofShard),
	}

}
