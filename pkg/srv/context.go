package srv

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	"github.com/ShatteredRealms/character-service/pkg/service"
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	commonrepo "github.com/ShatteredRealms/go-common-service/pkg/repository"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CharacterContext struct {
	*commonsrv.Context

	CharacterBusWriter characterbus.BusWriter

	CharacterService service.CharacterService
	InventoryService service.InventoryService

	DimensionService dimensionbus.Service
}

func NewCharacterContext(ctx context.Context, cfg *config.CharacterConfig) (*CharacterContext, error) {
	characterCtx := &CharacterContext{
		Context:            commonsrv.NewContext(&cfg.BaseConfig, config.ServiceName),
		CharacterBusWriter: bus.NewKafkaMessageBusWriter(config.GlobalConfig.BaseConfig.Kafka, characterbus.Message{}),
	}
	ctx, span := characterCtx.Tracer.Start(ctx, "context.character.new")
	defer span.End()

	pg, err := commonrepo.ConnectDB(ctx, cconfig.DBPoolConfig{Master: config.GlobalConfig.Postgres}, config.GlobalConfig.Redis)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	shouldMigrate := cfg.Mode != cconfig.ModeProduction
	migrater, err := commonrepo.NewPgxMigrater(
		ctx,
		config.GlobalConfig.Postgres.PostgresDSN(),
		config.GlobalConfig.MigrationPath,
		shouldMigrate,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres migrater: %w", err)
	}

	characterCtx.CharacterService = service.NewCharacterService(
		repository.NewPgxCharacterRepository(migrater),
	)
	characterCtx.DimensionService = dimensionbus.NewService(
		dimensionbus.NewPostgresRepository(pg),
		bus.NewKafkaMessageBusReader(config.GlobalConfig.Kafka, config.ServiceName, dimensionbus.Message{}),
	)
	characterCtx.DimensionService.StartProcessing(ctx)

	characterCtx.InventoryService = service.NewInventoryService(
		repository.NewPgxInventoryRepository(migrater),
	)

	return characterCtx, nil
}

func (c *CharacterContext) getDimension(ctx context.Context, dimensionId string) (*dimensionbus.Dimension, error) {
	dimension, err := c.DimensionService.GetDimensionById(ctx, dimensionId)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrDimensionLookup, err)
		return nil, status.Error(codes.Internal, ErrDimensionLookup.Error())
	}
	if dimension == nil {
		return nil, status.Error(codes.InvalidArgument, ErrDimensionNotExist.Error())
	}
	return dimension, nil
}

func (c *CharacterContext) getCharacterById(ctx context.Context, characterId string) (*character.Character, error) {
	id, err := uuid.Parse(characterId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrCharacterIdInvalid.Error())
	}

	character, err := c.CharacterService.GetCharacterById(ctx, &id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}

	if character == nil {
		return nil, status.Error(codes.NotFound, ErrCharacterDoesNotExist.Error())
	}
	return character, nil
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
		chars, _, err := c.CharacterService.GetCharacters(ctx)
		if err != nil {
			return fmt.Errorf("get characters: %w", err)
		}
		deletedChars, _, err := c.CharacterService.GetDeletedCharacters(ctx)
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
