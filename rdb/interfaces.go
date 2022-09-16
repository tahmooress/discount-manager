package rdb

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetAll(ctx context.Context, prefix string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Flush(ctx context.Context) error
	Close() error
}
