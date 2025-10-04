package repository

import (
	"database/sql"
	"time"

	sqlc "github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/gen"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
)

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
