package srv

import (
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	commonconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

type CharacterContext struct {
	*commonsrv.Context

	CharacterCreatedBus bus.MessageBus[bus.CharacterCreatedMessage]
}

func NewCharacterContext(baseConfig *commonconfig.BaseConfig, serviceName string) *CharacterContext {
	cfg := &CharacterContext{
		Context:             commonsrv.NewContext(baseConfig, serviceName),
		CharacterCreatedBus: bus.NewKafkaMessageBus(baseConfig.Kafka, serviceName, bus.CharacterCreatedMessage{}),
	}
	return cfg
}
