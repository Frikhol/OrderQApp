package config

type Config struct {
	ServiceName string   `envconfig:"SERVICE_NAME" required:"true"`
	Version     string   `envconfig:"VERSION" required:"true"`
	GRPCPort    string   `envconfig:"GRPC_PORT" default:"9000"`
	LogLevel    string   `envconfig:"LOG_LEVEL" default:"debug"`
	Postgres    Postgres `envconfig:"POSTGRES" required:"true"`
	RabbitMQ    RabbitMQ `envconfig:"RABBITMQ" required:"true"`
}

type Postgres struct {
	Host     string `envconfig:"HOST" default:"order_postgres"`
	Port     string `envconfig:"PORT" default:"5432"`
	Username string `envconfig:"USERNAME" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"postgres"`
	Database string `envconfig:"DATABASE" default:"order_db"`
	Secret   string `envconfig:"SECRET" default:"secret"`
}

type RabbitMQ struct {
	URL string `envconfig:"URL" default:"amqp://guest:guest@rabbitmq:5672/"`
}
