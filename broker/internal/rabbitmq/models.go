package rabbitmq

import (
	"time"

	"github.com/tahmooress/discount-manager/logger"
)

type Config struct {
	Addr              string
	ExchangeName      string
	ExchangeType      string
	Queue             string
	RouteKey          string
	ConsumerTag       string
	PrefetchCount     int64
	PrefetchSize      int64
	ReconnectInterval time.Duration
	ReInitInterval    time.Duration
	// in case of nil, deafault logger will use.
	logger logger.Logger
}

type HealthState struct {
	Er      error
	Status  bool
	Message string
}
