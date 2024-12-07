package service

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	"github.com/google/uuid"
)

var (
	ErrCharacter = errors.New("character service")
)

type CharacterService interface {
	GetCharacters(ctx context.Context) (*character.Models, error)
	GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Models, error)

	GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Model, error)

	CreateCharacter(ctx context.Context, ownerId, name, gender, realm string, dimension *dimensionbus.Dimension) (*character.Model, error)

	DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Model, error)

	EditCharacter(ctx context.Context, newCharacter *character.Model) (*character.Model, error)

	AddCharacterPlaytime(ctx context.Context, char *character.Model, seconds uint64) (*character.Model, error)
}

type characterService struct {
	repo repository.CharacterRepository
}

func NewCharacterService(repo repository.CharacterRepository) CharacterService {
	return &characterService{repo: repo}
}

// AddCharacterPlaytime implements CharacterService.
func (c *characterService) AddCharacterPlaytime(ctx context.Context, character *character.Model, seconds uint64) (*character.Model, error) {
	character.PlayTime += seconds
	return c.repo.UpdateCharacter(ctx, character)
}

// CreateCharacter implements CharacterService.
func (c *characterService) CreateCharacter(ctx context.Context, ownerId string, name string, gender string, realm string, dimension *dimensionbus.Dimension) (*character.Model, error) {
	character := &character.Model{
		Name:      name,
		OwnerId:   ownerId,
		Gender:    game.Gender(gender),
		Realm:     game.Realm(realm),
		Dimension: dimension,
	}

	err := character.Validate()
	if err != nil {
		return nil, err
	}

	return c.repo.CreateCharacter(ctx, character)
}

// DeleteCharacter implements CharacterService.
func (c *characterService) DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Model, error) {
	return c.repo.DeleteCharacter(ctx, characterId)
}

// EditCharacter implements CharacterService.
func (c *characterService) EditCharacter(ctx context.Context, character *character.Model) (*character.Model, error) {
	err := character.Validate()
	if err != nil {
		return nil, err
	}

	return c.repo.UpdateCharacter(ctx, character)
}

// GetCharacterById implements CharacterService.
func (c *characterService) GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Model, error) {
	return c.repo.GetCharacterById(ctx, characterId)
}

// GetCharacters implements CharacterService.
func (c *characterService) GetCharacters(ctx context.Context) (*character.Models, error) {
	return c.repo.GetCharacters(ctx)
}

// GetCharactersByOwner implements CharacterService.
func (c *characterService) GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Models, error) {
	return c.repo.GetCharactersByOwner(ctx, ownerId)
}
