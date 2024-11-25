package srv

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *characterServiceServer) getCharacterById(ctx context.Context, characterId string) (*character.Character, error) {
	character, err := c.Context.CharacterService.GetCharacterById(ctx, characterId)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %w", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}
	if character == nil {
		return nil, status.Error(codes.NotFound, ErrCharacterDoesNotExist.Error())
	}

	return character, nil
}
