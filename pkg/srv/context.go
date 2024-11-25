package srv

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	"github.com/ShatteredRealms/character-service/pkg/service"
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	commonrepo "github.com/ShatteredRealms/go-common-service/pkg/repository"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

type CharacterContext struct {
	*commonsrv.Context

	CharacterCreatedBus bus.MessageBus[bus.CharacterCreatedMessage]
	CharacterService    service.CharacterService
	DimensionService    service.DimensionService
}

func NewCharacterContext(ctx context.Context, cfg *config.CharacterConfig, serviceName string) (*CharacterContext, error) {
	characterCtx := &CharacterContext{
		Context:             commonsrv.NewContext(&cfg.BaseConfig, serviceName),
		CharacterCreatedBus: bus.NewKafkaMessageBus(cfg.Kafka, serviceName, bus.CharacterCreatedMessage{}),
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
	characterCtx.DimensionService = service.NewDimensionService(
		repository.NewPostgresDimensionRepository(pg),
	)

	return characterCtx, nil
}
