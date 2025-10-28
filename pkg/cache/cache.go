package cache

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/repository"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
)

type Cache struct {
	lru *lru.Cache[string, *model.Order] // {uid: Order}
}

// NewCache creates new cache size of N
func NewCache(size int, repo repository.OrderRepository) (*Cache, error) {

	// create empty map
	cache, err := lru.New[string, *model.Order](size)
	if err != nil {
		return nil, err
	}
	c := &Cache{lru: cache}

	ctx := context.Background()

	// get last N orders
	lastOrders, err := repo.GetLastOrders(ctx, size)
	if err != nil {
		return nil, err
	}

	for _, order := range lastOrders {

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
