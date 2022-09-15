package configs

import "os"

type AppConfigs struct {
	// HTTP server config
	HTTPIP           string
	HTTPPort         string
	HTTPReadTimeout  string
	HTTPWriteTimeout string

	// logger config
	LogLevel    string
	LogFilePath string

	// database config
	DatabaseDriver      string
	DatabaseName        string
	DatabaseUser        string
	DatabasePassword    string
	DatabaseHost        string
	DatabasePort        string
	DatabaseMaxPageSize string
	DatabaseSSLMode     string

	// rabbitMQ config
	RabbitAdr     string
	Exchange      string
	ExchangeType  string
	Queue         string
	RouteKey      string
	ConsumerTag   string
	PrefetchCount string
	PrefetchSize  string

	SentryEnv        string
	SentryDSN        string
	SentrySampleRate string
}

func Load() AppConfigs {
	return AppConfigs{
		HTTPIP:           os.Getenv("HTTP_IP"),
		HTTPPort:         os.Getenv("HTTP_PORT"),
		HTTPReadTimeout:  os.Getenv("HTTP_READ_TIMEOUT"),
		HTTPWriteTimeout: os.Getenv("HTTP_WRITE_TIMEOUT"),

		LogFilePath: os.Getenv("APP_LOG_PATH"),
		LogLevel:    os.Getenv("LOG_LEVEL"),

		DatabaseDriver:      os.Getenv("DATABASE_DRIVER"),
		DatabaseName:        os.Getenv("DATABASE_NAME"),
		DatabaseUser:        os.Getenv("DATABASE_USER"),
		DatabasePassword:    os.Getenv("DATABASE_PASS"),
		DatabaseHost:        os.Getenv("DATABASE_HOST"),
		DatabasePort:        os.Getenv("DATABASE_PORT"),
		DatabaseMaxPageSize: os.Getenv("DATABASE_MAX_PAGE_SIZE"),
		DatabaseSSLMode:     os.Getenv("DATABASE_SSLMODE"),

		RabbitAdr:     os.Getenv("RABBITMQ_ADDRESS"),
		Exchange:      os.Getenv("RABBITMQ_EXCHANGE"),
		ExchangeType:  os.Getenv("RABBITMQ_EXCHANGE_TYPE"),
		Queue:         os.Getenv("RABBITMQ_QUEUE"),
		RouteKey:      os.Getenv("RABBITMQ_ROUTING_KEY"),
		ConsumerTag:   os.Getenv("RABBITMQ_CONSUMER_TAG"),
		PrefetchCount: os.Getenv("RABBITMQ_PREFETCH_COUNT"),
		PrefetchSize:  os.Getenv("RABBITMQ_PREFETCH_SIZE"),

		SentryEnv:        os.Getenv("SENTRY_ENVIRONMENT"),
		SentryDSN:        os.Getenv("SENTRY_DSN"),
		SentrySampleRate: os.Getenv("SENTRY_TRACES_SAMPLE_RATE"),
	}
}
