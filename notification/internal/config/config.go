package config

type Config struct {
	ServiceName string   `envconfig:"SERVICE_NAME" required:"true"`
	Version     string   `envconfig:"VERSION" required:"true"`
	Port        string   `envconfig:"PORT" default:"8080"`
	LogLevel    string   `envconfig:"LOG_LEVEL" default:"debug"`
	RabbitMQ    RabbitMQ `envconfig:"RABBITMQ" required:"true"`
}

type RabbitMQ struct {
	URL string `envconfig:"URL" default:"amqp://guest:guest@rabbitmq:5672/"`
}
