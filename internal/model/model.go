package model

import "time"

type Order struct {
	OrderUID    string  `json:"order_uid" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Entry       *string `json:"entry" validate:"required"`

	Delivery Delivery `json:"delivery" validate:"required"`
	Payment  Payment  `json:"payment" validate:"required"`
	Items    []*Item  `json:"items" validate:"required,dive"`

	Locale            *string    `json:"locale" validate:"required"`
	InternalSignature *string    `json:"internal_signature"`
	CustomerID        *string    `json:"customer_id" validate:"required"`
	DeliveryService   *string    `json:"delivery_service" validate:"required"`
	ShardKey          *string    `json:"shardkey" validate:"required"`
	SmID              *int64     `json:"sm_id" validate:"required"`
	DateCreated       *time.Time `json:"date_created,omitempty" validate:"required,notfuture"`
	OofShard          *string    `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	OrderUID string `json:"-"`
	Name     string `json:"name"     validate:"required,name"`
	Phone    string `json:"phone"    validate:"required,phone"`
	Zip      string `json:"zip"      validate:"required,zip"`
	City     string `json:"city"     validate:"required,city"`
	Address  string `json:"address"  validate:"required,address"`
	Region   string `json:"region"   validate:"required,region"`
	Email    string `json:"email"    validate:"required,email"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    *string `json:"request_id"`
	Currency     *string `json:"currency" validate:"required"`
	Provider     *string `json:"provider" validate:"required"`
	Amount       *int64  `json:"amount" validate:"required"`
	PaymentDT    *int64  `json:"payment_dt" validate:"required"`
	Bank         *string `json:"bank" validate:"required"`
	DeliveryCost *int64  `json:"delivery_cost" validate:"required"`
	GoodsTotal   *int64  `json:"goods_total" validate:"required"`
	CustomFee    *int64  `json:"custom_fee" validate:"required"`
}

type Item struct {
	ID          int     `json:"-"`
	OrderUID    string  `json:"-"`
	ChrtID      *int64  `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       *int64  `json:"price" validate:"required"`
	Rid         *string `json:"rid" validate:"required"`
	Name        *string `json:"name" validate:"required"`
	Sale        *int64  `json:"sale" validate:"required"`
	Size        *string `json:"size" validate:"required"`
	TotalPrice  *int64  `json:"total_price" validate:"required"`
	NmID        *int64  `json:"nm_id" validate:"required"`
	Brand       *string `json:"brand" validate:"required"`
	Status      *int64  `json:"status" validate:"required"`
}
