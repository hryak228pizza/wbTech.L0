package consumer

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/hryak228pizza/wbTech.L0/internal/config"
	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
	"github.com/hryak228pizza/wbTech.L0/pkg/cache"
	"github.com/hryak228pizza/wbTech.L0/pkg/validation"
	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// runs a consumer for reading all incoming orders with kafka
func Consumer(ctx context.Context, cfg *config.Config, cache *cache.Cache, db *sql.DB) {

	// init reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.KafkaBroker},
		Topic:          cfg.KafkaTopic,
		GroupID:        cfg.KafkaGroup,
		CommitInterval: 0,
	})
	defer r.Close()

	logger.L().Info("consumer subscribe",
		zap.String("broker", r.Config().Brokers[0]),
		zap.String("topic", r.Config().Topic),
		zap.String("group", r.Config().GroupID),
	)

	// init validator
	validate := validation.NewValidator()

	for {
		// check context cancel
		select {
		case <-ctx.Done():
			logger.L().Info("consumer context canceled, exiting loop")
			return
		default:
		}

		innerCtx := context.Background()

		// fetch message
		m, err := r.FetchMessage(innerCtx)
		if err != nil {
			if ctx.Err() != nil {
				logger.L().Info("fetch aborted by context", zap.Error(err))
				return
			}
			logger.L().Error("kafka fetch failed",
				zap.String("topic", r.Config().Topic),
				zap.String("group", r.Config().GroupID),
				zap.String("error", err.Error()),
			)
			continue
		}

		// json deserialization
		var order model.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			logger.L().Error("json parsing failed",
				zap.String("error", err.Error()),
			)
			continue
		}

		// validate order data
		if err := validate.ValidateOrder(&order); err != nil {
			if verrs, ok := err.(validator.ValidationErrors); ok {
				for _, verr := range verrs {
					logger.L().Error("validation error",
						zap.String("order_uid", order.OrderUID),
						zap.String("field", verr.Namespace()),
						zap.String("tag", verr.Tag()),
						zap.String("param", verr.Param()),
						zap.Any("value", verr.Value()),
					)
				}
			} else {
				logger.L().Error("validation failed",
					zap.String("order_uid", order.OrderUID),
					zap.Error(err),
				)
			}
			logger.L().Warn("invalid order skipped",
				zap.String("order_uid", order.OrderUID),
			)
			continue
		}

		// save order in database
		if err := saveOrder(db, &order, cache); err != nil {
			logger.L().Error("db writing failed",
				zap.String("error", err.Error()),
			)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			logger.L().Error("commit failed",
				zap.Error(err),
			)
		} else {
			logger.L().Info("order saved and committed",
				zap.String("order_id", order.OrderUID),
			)
		}
	}
}

// saves order in DB
func saveOrder(db *sql.DB, o *model.Order, cache *cache.Cache) error {

	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// inserts:
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

	// save order into cache
	cache.SetOrder(o)

	return tx.Commit()
}
