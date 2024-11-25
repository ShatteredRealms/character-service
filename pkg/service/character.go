package service

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/repository"
)

var (
	ErrCharacter = errors.New("character service")
)

type CharacterService interface {
	GetCharacters(ctx context.Context) (*character.Characters, error)
	GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Characters, error)

	GetCharacterById(ctx context.Context, characterId string) (*character.Character, error)

	CreateCharacter(ctx context.Context, ownerId, name, gender, realm, dimension string) (*character.Character, error)

	DeleteCharacter(ctx context.Context, characterId string) (*character.Character, error)

	EditCharacter(ctx context.Context, character *character.Character) (*character.Character, error)

	AddCharacterPlaytime(ctx context.Context, character *character.Character, seconds uint64) (*character.Character, error)
}

type characterService struct {
	repo repository.CharacterRepository
}

func NewCharacterService(repo repository.CharacterRepository) CharacterService {
	return &characterService{repo: repo}
}

// AddCharacterPlaytime implements CharacterService.
func (c *characterService) AddCharacterPlaytime(ctx context.Context, character *character.Character, seconds uint64) (*character.Character, error) {
	character.PlayTime += seconds
	return c.repo.UpdateCharacter(ctx, character)
}

// CreateCharacter implements CharacterService.
func (c *characterService) CreateCharacter(ctx context.Context, ownerId string, name string, gender string, realm string, dimension string) (*character.Character, error) {
	character := &character.Character{
		Name:      name,
		Gender:    game.Gender(gender),
		Realm:     game.Realm(realm),
		Dimension: &game.Dimension{Id: dimension},
	}

	err := character.Validate()
	if err != nil {
		return nil, err
	}

	return c.repo.CreateCharacter(ctx, character)
}

// DeleteCharacter implements CharacterService.
func (c *characterService) DeleteCharacter(ctx context.Context, characterId string) (*character.Character, error) {
	return c.repo.DeleteCharacter(ctx, characterId)
}

// EditCharacter implements CharacterService.
func (c *characterService) EditCharacter(ctx context.Context, character *character.Character) (*character.Character, error) {
	err := character.Validate()
	if err != nil {
		return nil, err
	}

	return nil, ErrCharacter
}

// GetCharacterById implements CharacterService.
func (c *characterService) GetCharacterById(ctx context.Context, characterId string) (*character.Character, error) {
	return c.repo.GetCharacterById(ctx, characterId)
}

// GetCharacters implements CharacterService.
func (c *characterService) GetCharacters(ctx context.Context) (*character.Characters, error) {
	return c.repo.GetCharacters(ctx)
}

// GetCharactersByOwner implements CharacterService.
func (c *characterService) GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Characters, error) {
	return c.repo.GetCharactersByOwner(ctx, ownerId)
}
