package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hryak228pizza/wbTech.L0/internal/config"
	"github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/repository"
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

		// process message
		//innerCtx := context.Background()
		if err := processMessage(ctx, r, db, cache, validate); err != nil {
			logger.L().Error("message processing failed",
				zap.String("error", err.Error()),
			)
		}
	}
}

func processMessage(ctx context.Context, r *kafka.Reader, db *sql.DB, cache *cache.Cache, validate *validation.Validate) error {

	repo := repository.NewOrderRepo(db)

	// fetch message
	m, err := r.FetchMessage(ctx)
	if err != nil {
		if ctx.Err() != nil {
			logger.L().Info("fetch aborted by context", zap.Error(err))
			return err
		}
		logger.L().Error("kafka fetch failed",
			zap.String("error", err.Error()),
		)
		time.Sleep(2 * time.Second)
		return err
	}

	// json deserialization
	var order model.Order
	if err := json.Unmarshal(m.Value, &order); err != nil {
		logger.L().Error("json parsing failed",
			zap.String("error", err.Error()),
		)
		return err
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
		return err
	}

	// save order in database
	if err := repo.Save(ctx, &order); err != nil {
		logger.L().Error("db writing failed",
			zap.String("order_uid", order.OrderUID),
			zap.String("error", err.Error()),
		)
		return err
	}

	// save order into cache
	cache.SetOrder(&order)

	// commit order
	commitContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.CommitMessages(commitContext, m); err != nil {
		logger.L().Error("commit failed",
			zap.Error(err),
		)
	} else {
		logger.L().Info("order saved and committed",
			zap.String("order_id", order.OrderUID),
		)
	}

	return nil
}
