package service

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/gender"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/profession"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/realm"
	"github.com/google/uuid"
)

var (
	ErrCharacter = errors.New("character service")
)

type CharacterService interface {
	GetCharacters(ctx context.Context) (character.Characters, int, error)
	GetDeletedCharacters(ctx context.Context) (character.Characters, int, error)
	GetCharactersByOwner(ctx context.Context, ownerId string) (character.Characters, int, error)

	GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Character, error)

	CreateCharacter(ctx context.Context, ownerId, name, gender, realm string, profession string, dimensionId *uuid.UUID) (*character.Character, error)

	DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error)

	EditCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error)

	AddCharacterPlaytime(ctx context.Context, char *character.Character, seconds int32) (*character.Character, error)
}

type characterService struct {
	repo repository.CharacterRepository
}

func NewCharacterService(repo repository.CharacterRepository) CharacterService {
	return &characterService{repo: repo}
}

// AddCharacterPlaytime implements CharacterService.
func (c *characterService) AddCharacterPlaytime(ctx context.Context, character *character.Character, seconds int32) (*character.Character, error) {
	character.PlayTime += seconds
	return c.repo.UpdateCharacter(ctx, character)
}

// CreateCharacter implements CharacterService.
func (c *characterService) CreateCharacter(
	ctx context.Context,
	ownerId string,
	name string,
	genderStr string,
	realmStr string,
	professionStr string,
	dimensionId *uuid.UUID,
) (*character.Character, error) {
	ownerIdUUID, err := uuid.Parse(ownerId)
	if err != nil {
		return nil, err
	}

	_, count, err := c.GetCharactersByOwner(ctx, ownerId)
	if err != nil {
		return nil, err
	}
	if count >= 10 {
		return nil, errors.New("too many characters")
	}

	character := &character.Character{
		Name:        name,
		OwnerId:     ownerIdUUID,
		Gender:      gender.Gender(genderStr),
		Realm:       realm.Realm(realmStr),
		Profession:  profession.Profession(professionStr),
		DimensionId: *dimensionId,
	}

	err = character.Validate()
	if err != nil {
		return nil, err
	}

	return c.repo.CreateCharacter(ctx, character)
}

// DeleteCharacter implements CharacterService.
func (c *characterService) DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	return c.repo.DeleteCharacter(ctx, characterId)
}

// EditCharacter implements CharacterService.
func (c *characterService) EditCharacter(ctx context.Context, character *character.Character) (*character.Character, error) {
	err := character.Validate()
	if err != nil {
		return nil, err
	}

	return c.repo.UpdateCharacter(ctx, character)
}

// GetCharacterById implements CharacterService.
func (c *characterService) GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	return c.repo.GetCharacter(ctx, map[string]any{"id": characterId})
}

// GetCharacters implements CharacterService.
func (c *characterService) GetCharacters(ctx context.Context) (character.Characters, int, error) {
	return c.repo.GetCharacters(ctx, nil, nil, false)
}

// GetDeletedCharacters implements CharacterService.
func (c *characterService) GetDeletedCharacters(ctx context.Context) (character.Characters, int, error) {
	return c.repo.GetCharacters(ctx, nil, nil, true)
}

// GetCharactersByOwner implements CharacterService.
func (c *characterService) GetCharactersByOwner(ctx context.Context, ownerId string) (character.Characters, int, error) {
	id, err := uuid.Parse(ownerId)
	if err != nil {
		return nil, -1, ErrInvalidOwnerId
	}
	return c.repo.GetCharacters(ctx, map[string]any{"owner_id": id}, nil, false)
}
