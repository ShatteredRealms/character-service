package config

import (
	"context"

	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/sirupsen/logrus"
)

var (
	Version         = "v1.0.0"
	ServiceName     = "CharacterService"
	GlobalConfig    *CharacterConfig
	GlobalConfigErr error
)

func init() {
	GlobalConfig, GlobalConfigErr = NewCharacterConfig(context.Background())
}

type CharacterConfig struct {
	cconfig.BaseConfig `yaml:",inline" mapstructure:",squash"`
	Postgres           cconfig.DBConfig     `yaml:"postgres"`
	Redis              cconfig.DBPoolConfig `yaml:"redis"`
}

func NewCharacterConfig(ctx context.Context) (*CharacterConfig, error) {
	config := &CharacterConfig{
		BaseConfig: cconfig.BaseConfig{
			Server: cconfig.ServerAddress{
				Host: "localhost",
				Port: "8081",
			},
			Keycloak: cconfig.KeycloakConfig{
				BaseURL:      "http://localhost:8080",
				Realm:        "default",
				Id:           "738a426a-da91-4b16-b5fc-92d63a22eb76",
				ClientId:     "sro-character-service",
				ClientSecret: "**********",
			},
			Mode:                "local",
			LogLevel:            logrus.InfoLevel,
			OpenTelemtryAddress: "localhost:4317",
			Kafka: cconfig.ServerAddresses{
				{
					Host: "localhost",
					Port: "29092",
				},
			},
			MigrationPath: "migrations",
		},
		Postgres: cconfig.DBConfig{
			ServerAddress: cconfig.ServerAddress{
				Host: "localhost",
				Port: "5432",
			},
			Name:     "character_service",
			Username: "postgres",
			Password: "password",
		},
		Redis: cconfig.DBPoolConfig{
			Master: cconfig.DBConfig{
				ServerAddress: cconfig.ServerAddress{
					Host: "localhost",
					Port: "7000",
				},
			},
		},
	}

	err := cconfig.BindConfigEnvs(ctx, "sro-character", config)
	return config, err
}
