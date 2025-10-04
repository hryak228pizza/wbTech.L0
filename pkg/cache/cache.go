package cache

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	sqlc "github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/gen"
	"github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/repository"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
)

type Cache struct {
	lru *lru.Cache[string, *model.Order] // {uid: Order}
}

// NewCache creates new cache size of N
func NewCache(size int, dbQueries *sqlc.Queries) (*Cache, error) {

	// create empty map
	cache, err := lru.New[string, *model.Order](size)
	if err != nil {
		return nil, err
	}
	c := &Cache{lru: cache}

	ctx := context.Background()

	// get last N orders
	lastOrders, err := dbQueries.GetLastOrders(ctx, int32(size))
	if err != nil {
		return nil, err
	}

	for _, elem := range lastOrders {

		// parse delivery
		delivery, err := dbQueries.GetDeliveryByOrderUID(ctx, elem.OrderUid)
		if err != nil {
			return nil, err
		}

		// parse payment
		payment, err := dbQueries.GetPaymentByTransaction(ctx, elem.OrderUid)
		if err != nil {
			return nil, err
		}

		// parse items
		items, err := dbQueries.GetItemsByTrackNumber(ctx, elem.TrackNumber)
		if err != nil {
			return nil, err
		}

		// create full order from piecies
		order := repository.MapToOrder(elem, delivery, payment, items)

		// save order to cache
		c.SetOrder(order)
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
