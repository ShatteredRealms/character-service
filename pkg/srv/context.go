package srv

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	"github.com/ShatteredRealms/character-service/pkg/service"
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	commonrepo "github.com/ShatteredRealms/go-common-service/pkg/repository"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

type CharacterContext struct {
	*commonsrv.Context

	CharacterBusWriter bus.MessageBusWriter[characterbus.Message]

	CharacterService service.CharacterService

	DimensionService dimensionbus.Service
}

func NewCharacterContext(ctx context.Context, cfg *config.CharacterConfig, serviceName string) (*CharacterContext, error) {
	characterCtx := &CharacterContext{
		Context:            commonsrv.NewContext(&cfg.BaseConfig, serviceName),
		CharacterBusWriter: bus.NewKafkaMessageBusWriter(cfg.Kafka, characterbus.Message{}),
	}
	ctx, span := characterCtx.Tracer.Start(ctx, "context.character.new")
	defer span.End()

	pg, err := commonrepo.ConnectDB(ctx, cfg.Postgres, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	characterCtx.CharacterService = service.NewCharacterService(
		repository.NewPostgresCharacterRepository(pg),
	)
	characterCtx.DimensionService = dimensionbus.NewService(
		dimensionbus.NewPostgresRepository(pg),
		bus.NewKafkaMessageBusReader(cfg.Kafka, serviceName, dimensionbus.Message{}),
	)
	characterCtx.DimensionService.StartProcessing(ctx)

	return characterCtx, nil
}
