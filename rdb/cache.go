package rdb

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	errEmptyAddress = errors.New("address should not be empty")
)

const defaultDB = 0

type cache struct {
	client *redis.Client
}

func New(host, port, db string) (Cache, error) {
	if host == "" || port == "" {
		return nil, errEmptyAddress
	}

	n, err := getDB(db)
	if err != nil {
		return nil, err
	}

	return &cache{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", host, port),
			DB:   n,
		}),
	}, nil
}

func getDB(db string) (int, error) {
	var n int64

	if db == "" {
		n = defaultDB
	} else {
		i, err := strconv.ParseInt(db, 10, 64)
		if err != nil {
			return 0, err
		}

		if i <= 15 {
			n = i
		}
	}

	return int(n), nil
}
