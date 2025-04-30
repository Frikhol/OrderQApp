package config

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`
	Version     string `envconfig:"VERSION" required:"true"`
	GRPCPort    string `envconfig:"GRPC_PORT" default:"9000"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
	Host        string `envconfig:"HOST" default:"postgres"`
	Port        string `envconfig:"PORT" default:"5432"`
	Username    string `envconfig:"USERNAME" default:"postgres"`
	Password    string `envconfig:"PASSWORD" default:"postgres"`
	Database    string `envconfig:"DATABASE" default:"auth_db"`
	Secret      string `envconfig:"SECRET" default:"secret"`
}
