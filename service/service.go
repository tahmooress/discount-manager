package service

import (
	"encoding/json"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tahmooress/discount-manager/broker/rabbitmq"
	"github.com/tahmooress/discount-manager/configs"
	"github.com/tahmooress/discount-manager/entities"
	"github.com/tahmooress/discount-manager/logger"
	"github.com/tahmooress/discount-manager/pkg/wrapper"
	"github.com/tahmooress/discount-manager/rdb"
	"github.com/tahmooress/discount-manager/repository"
)

type service struct {
	db      repository.DB
	cachedb rdb.Cache

	queue    rabbitmq.Publisher
	notifier rabbitmq.Publisher
	consumer rabbitmq.Consumer

	activeCampaigns map[string]entities.Campaign
	mu              sync.RWMutex

	closer *wrapper.Closer
	logger logger.Logger
}

func New(cfg *configs.AppConfigs, logger logger.Logger) (u Usecases, err error) {
	c := new(wrapper.Closer)

	db, err := repository.New(cfg)
	if err != nil {
		return nil, err
	}

	c.Add(db)

	defer func() {
		if err != nil {
			_ = c.Close()
		}
	}()

	notifier, err := rabbitmq.NewPublisher(rabbitmq.Config{
		Host:         cfg.RabbitMQWalletHost,
		Port:         cfg.RabbitMQWalletPort,
		User:         cfg.RabbitMQWalletUser,
		Pass:         cfg.RabbitMQWalletPass,
		ExchangeName: cfg.RabbitMQWalletExchange,
		ExchangeType: cfg.RabbitMQWalletExchangeType,
		RouteKey:     cfg.RabbitMQWalletRoutingKey,
		Queue:        cfg.RabbitMQWalletQuee,
	}, logger)
	if err != nil {
		return nil, err
	}

	c.Add(notifier)

	queue, err := rabbitmq.NewPublisher(rabbitmq.Config{
		Host:         cfg.RabbitMQRedeemerHost,
		Port:         cfg.RabbitMQRedeemerPort,
		User:         cfg.RabbitMQWalletUser,
		Pass:         cfg.RabbitMQWalletPass,
		ExchangeName: cfg.RabbitMQRedeemerExchange,
		ExchangeType: cfg.RabbitMQRedeemerExchangeType,
		RouteKey:     cfg.RabbitMQRedeemerRouteKey,
		Queue:        cfg.RabbitMQRedeemerQueue,
	}, logger)
	if err != nil {
		return nil, err
	}

	c.Add(queue)

	cachedb, err := rdb.New(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		return nil, err
	}

	c.Add(cachedb)

	s := &service{
		db:       db,
		cachedb:  cachedb,
		queue:    queue,
		notifier: notifier,
		mu:       sync.RWMutex{},
		closer:   c,
		logger:   logger,
	}

	err = s.intiCampagins()
	if err != nil {
		return nil, err
	}

	err = s.initConsumer(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// register handler to consumer and start consuming.
func (s *service) initConsumer(cfg *configs.AppConfigs) error {
	handler := func(d amqp.Delivery) (action rabbitmq.Action) {
		var redeemer entities.Redeemer

		err := json.Unmarshal(d.Body, &redeemer)
		if err != nil {
			return rabbitmq.NackDiscard
		}

		err = s.ApplyVoucher(&redeemer)
		if err != nil {
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	}

	consumer, err := rabbitmq.NewConsumer(
		rabbitmq.Config{
			Host:         cfg.RabbitMQRedeemerHost,
			Port:         cfg.RabbitMQRedeemerPort,
			User:         cfg.RabbitMQRedeemerUser,
			Pass:         cfg.RabbitMQRedeemerPass,
			ExchangeName: cfg.RabbitMQRedeemerExchange,
			ExchangeType: cfg.RabbitMQRedeemerExchangeType,
			RouteKey:     cfg.RabbitMQRedeemerRouteKey,
			Queue:        cfg.RabbitMQRedeemerQueue,
		},
		handler,
		s.logger,
	)
	if err != nil {
		return err
	}

	s.consumer = consumer
	s.closer.Add(consumer)

	return nil
}

func (s *service) intiCampagins() error {
	campaigns, err := s.db.GetCampaignsByStatus(true)
	if err != nil {
		return err
	}

	s.activeCampaigns = make(map[string]entities.Campaign)

	for _, cmp := range campaigns {
		s.activeCampaigns[cmp.Name] = cmp
	}

	return nil
}

func (s *service) Close() error {
	return s.closer.Close()
}
