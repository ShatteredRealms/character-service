package repository

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/google/uuid"
)

type CharacterRepository interface {
	GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Model, error)

	GetCharacters(ctx context.Context) (*character.Models, error)
	GetCharactersByOwner(ctx context.Context, ownerId string) (*character.Models, error)

	CreateCharacter(ctx context.Context, newCharacter *character.Model) (*character.Model, error)
	UpdateCharacter(ctx context.Context, updatedCharacter *character.Model) (*character.Model, error)

	DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Model, error)
	DeleteCharactersByOwner(ctx context.Context, ownerId string) (*character.Models, error)
}
