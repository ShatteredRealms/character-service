package repository

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/google/uuid"
)

type CharacterRepository interface {
	// GetCharacters gets the of character that matches the provided filters. If more than one character matches the filters, an error is returned.
	GetCharacter(ctx context.Context, matchFilters map[string]interface{}) (*character.Character, error)

	// GetCharacters gets a list of characters that match the provided filters and returns the total number of characters that match the filters.
	// If no characters match the filters, an empty list is returned.
	// Deleted characters are only returned if deleted is true, and non-deleted characters if true.
	GetCharacters(ctx context.Context, matchFilters map[string]interface{}, queryFilter *pb.QueryFilters, deleted bool) (character.Characters, int, error)

	CreateCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error)
	UpdateCharacter(ctx context.Context, updatedCharacter *character.Character) (*character.Character, error)

	DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error)
	DeleteCharactersByOwner(ctx context.Context, ownerId *uuid.UUID) (character.Characters, error)
}
