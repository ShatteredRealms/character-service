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
	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
	commonrepo "github.com/ShatteredRealms/go-common-service/pkg/repository"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

type CharacterContext struct {
	*commonsrv.Context

	CharacterBusWriter characterbus.BusWriter

	CharacterService service.CharacterService

	DimensionService dimensionbus.Service
}

func NewCharacterContext(ctx context.Context) (*CharacterContext, error) {
	characterCtx := &CharacterContext{
		Context:            commonsrv.NewContext(&config.GlobalConfig.BaseConfig, config.ServiceName),
		CharacterBusWriter: bus.NewKafkaMessageBusWriter(config.GlobalConfig.BaseConfig.Kafka, characterbus.Message{}),
	}
	ctx, span := characterCtx.Tracer.Start(ctx, "context.character.new")
	defer span.End()

	pg, err := commonrepo.ConnectDB(ctx, cconfig.DBPoolConfig{Master: config.GlobalConfig.Postgres}, config.GlobalConfig.Redis)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	migrater, err := commonrepo.NewPgxMigrater(ctx, config.GlobalConfig.Postgres.PostgresDSN(), config.GlobalConfig.MigrationPath)
	if err != nil {
		return nil, fmt.Errorf("postgres migrater: %w", err)
	}

	repo := repository.NewPgxCharacterRepository(migrater)

	characterCtx.CharacterService = service.NewCharacterService(
		repo,
	)
	characterCtx.DimensionService = dimensionbus.NewService(
		dimensionbus.NewPostgresRepository(pg),
		bus.NewKafkaMessageBusReader(config.GlobalConfig.Kafka, config.ServiceName, dimensionbus.Message{}),
	)
	characterCtx.DimensionService.StartProcessing(ctx)

	return characterCtx, nil
}

func (c *CharacterContext) Close() {
	if c.DimensionService != nil {
		c.DimensionService.StopProcessing()
	}
}

func (c *CharacterContext) ResetCharacterBus() commonsrv.WriterResetCallback {
	return func(ctx context.Context) error {
		ctx, span := c.Context.Tracer.Start(ctx, "character.reset_character_bus")
		defer span.End()
		chars, err := c.CharacterService.GetCharacters(ctx)
		if err != nil {
			return fmt.Errorf("get characters: %w", err)
		}
		deletedChars, err := c.CharacterService.GetDeletedCharacters(ctx)
		if err != nil {
			return fmt.Errorf("get deleted characters: %w", err)
		}

		msgs := make([]characterbus.Message, len(chars)+len(deletedChars))
		for idx, char := range chars {
			msgs[idx] = characterbus.Message{
				Id:          char.Id,
				OwnerId:     char.OwnerId,
				DimensionId: char.DimensionId,
				MapId:       char.Location.WorldId,
				Deleted:     false,
			}
		}
		for idx, char := range deletedChars {
			msgs[idx+len(chars)] = characterbus.Message{
				Id:          char.Id,
				OwnerId:     char.OwnerId,
				DimensionId: char.DimensionId,
				MapId:       char.Location.WorldId,
				Deleted:     true,
			}
		}

		return c.CharacterBusWriter.PublishMany(ctx, msgs)
	}
}
