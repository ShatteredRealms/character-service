package repository

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/google/uuid"
)

type CharacterRepository interface {
	GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Character, error)

	GetCharacters(ctx context.Context) (*character.Characters, error)
	GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Characters, error)

	CreateCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error)
	UpdateCharacter(ctx context.Context, updatedCharacter *character.Character) (*character.Character, error)

	DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error)
	DeleteCharactersByOwner(ctx context.Context, ownerId string) (*character.Characters, error)
}
