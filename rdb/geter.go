package rdb

import (
	"context"
	"fmt"
)

func (c *cache) Close() error {
	return c.client.Close()
}

func (c *cache) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, key, value, 0).Err()
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *cache) GetAll(ctx context.Context, prefix string) ([]string, error) {
	itr := c.client.Scan(ctx, 0, fmt.Sprintf("%s:*", prefix), 0).Iterator()

	result := make([]string, 0)

	for itr.Next(ctx) {

		r, err := c.client.Get(ctx, itr.Val()).Result()
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	if err := itr.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *cache) Flush(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}
