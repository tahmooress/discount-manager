package configs

import "os"

type AppConfigs struct {
	// HTTP server config
	HTTPIP   string
	HTTPPort string

	// logger config
	LogLevel    string
	LogFilePath string

	// database config
	DatabaseDriver   string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	// rabbitMQ config
	RabbitMQRedeemerHost         string
	RabbitMQRedeemerPort         string
	RabbitMQRedeemerExchange     string
	RabbitMQRedeemerExchangeType string
	RabbitMQRedeemerQueue        string
	RabbitMQRedeemerRouteKey     string

	RabbitMQWalletHost         string
	RabbitMQWalletPort         string
	RabbitMQWalletExchange     string
	RabbitMQWalletExchangeType string
	RabbitMQWalletRoutingKey   string
	RabbitMQWalletQuee         string
}

func Load() *AppConfigs {
	return &AppConfigs{
		HTTPIP:   os.Getenv("HTTP_IP"),
		HTTPPort: os.Getenv("HTTP_PORT"),

		LogFilePath: os.Getenv("APP_LOG_PATH"),
		LogLevel:    os.Getenv("LOG_LEVEL"),

		DatabaseDriver:   os.Getenv("DATABASE_DRIVER"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASS"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseSSLMode:  os.Getenv("DATABASE_SSLMODE"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),

		RabbitMQRedeemerHost:         os.Getenv("RABBITMQ_REDEEMER_HOST"),
		RabbitMQRedeemerPort:         os.Getenv("RABBITMQ_REDEEMER_PORT"),
		RabbitMQRedeemerExchange:     os.Getenv("RABBITMQ_REDEEMER_EXCHANGE"),
		RabbitMQRedeemerExchangeType: os.Getenv("RABBITMQ_REDEEMER_EXCHANGE_TYPE"),
		RabbitMQRedeemerQueue:        os.Getenv("RABBITMQ_REDEEMER_QUEUE"),
		RabbitMQRedeemerRouteKey:     os.Getenv("RABBITMQ_REDEEMER_ROUTING_KEY"),

		RabbitMQWalletHost:         os.Getenv("RABBITMQ_WALLET_HOST"),
		RabbitMQWalletPort:         os.Getenv("RABBITMQ_REDEEMER_PORT"),
		RabbitMQWalletExchange:     os.Getenv("RABBITMQ_WALLET_EXCHANGE"),
		RabbitMQWalletExchangeType: os.Getenv("RABBITMQ_WALLET_EXCHANGE_TYPE"),
		RabbitMQWalletQuee:         os.Getenv("RABBITMQ_WALLET_QUEUE"),
		RabbitMQWalletRoutingKey:   os.Getenv("RABBITMQ_WALLET_ROUTING_KEY"),
	}
}
