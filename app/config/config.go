package config

import "github.com/caarlos0/env/v9"

type DBConfig struct {
	DBname string `env:"POSTGRES_DB" envDefault:"servicedb"`
	User   string `env:"POSTGRES_SERVICE_USER" envDefault:"serviceuser"`
	Pass   string `env:"POSTGRES_SERVICE_PASSWORD,unset" envDefault:"servicepassword"`
	Host   string `envDefault:"127.0.0.1"`
	Port   string `envDefault:"5432"`
}

type NATSConfig struct {
	ClusterID string `env:"NATS_CLUSTER_ID" envDefault:"test-cluster"`
	ClientID  string `env:"NATS_CLIENT_ID" envDefault:"test-client"`
}

type Config struct {
	DB   DBConfig
	NATS NATSConfig
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)

	return &cfg, err
}
