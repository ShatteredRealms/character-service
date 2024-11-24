package config

import (
	"context"

	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/sirupsen/logrus"
)

var (
	Version     = "v1.0.0"
	ServiceName = "CharacterService"
)

type CharacterConfig struct {
	cconfig.BaseConfig `yaml:",inline" characterstructure:",squash"`
	Postgres           cconfig.DBPoolConfig `yaml:"postgres"`
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
				Id:           "7b575e9b-c687-4cdc-b210-67c59b5f380f",
				ClientId:     "sro-character-service",
				ClientSecret: "**********",
			},
			Mode:                "local",
			LogLevel:            logrus.DebugLevel,
			OpenTelemtryAddress: "localhost:4317",
		},
		Postgres: cconfig.DBPoolConfig{
			Master: cconfig.DBConfig{
				ServerAddress: cconfig.ServerAddress{},
				Name:          "character-service",
				Username:      "postgres",
				Password:      "password",
			},
		},
	}

	err := cconfig.BindConfigEnvs(ctx, "sro-character", config)
	return config, err
}
