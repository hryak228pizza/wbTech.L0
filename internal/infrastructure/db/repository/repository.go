package repository

import (
	"context"
	"database/sql"

	sqlc "github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/gen"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
)

type OrderRepository interface {
	// Save saves order
	Save(ctx context.Context, o *model.Order) error

	// GetByUID returns full order by order_uid
	GetByUID(ctx context.Context, uid string) (*model.Order, error)

	// GetLastOrders returns last N entries
	GetLastOrders(ctx context.Context, limit int) ([]*model.Order, error)
}

type orderRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewOrderRepository(db *sql.DB, q *sqlc.Queries) OrderRepository {
	return &orderRepository{
		db:      db,
		queries: q,
	}
}
