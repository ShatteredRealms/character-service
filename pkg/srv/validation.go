package srv

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/auth"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *characterServiceServer) validateUserPermissions(ctx context.Context, ownerId string, selfRole, otherRole *gocloak.Role) error {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return commonsrv.ErrPermissionDenied
	}
	if !claims.HasResourceRole(selfRole, c.Context.Config.Keycloak.ClientId) {
		return commonsrv.ErrPermissionDenied
	}
	if claims.Subject != ownerId && !claims.HasResourceRole(otherRole, c.Context.Config.Keycloak.ClientId) {
		return commonsrv.ErrPermissionDenied
	}
	return nil
}

func (c *characterServiceServer) validateCharacterPermissions(ctx context.Context, characterId string, selfRole, otherRole *gocloak.Role) (*character.Character, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, commonsrv.ErrPermissionDenied
	}

	if !claims.HasResourceRole(selfRole, c.Context.Config.Keycloak.ClientId) {
		return nil, commonsrv.ErrPermissionDenied
	}

	id, err := uuid.Parse(characterId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrCharacterIdInvalid.Error())
	}

	character, err := c.Context.CharacterService.GetCharacterById(ctx, &id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}

	if character == nil {
		return nil, status.Error(codes.NotFound, ErrCharacterDoesNotExist.Error())
	}

	if claims.Subject != character.OwnerId.String() && !claims.HasResourceRole(otherRole, c.Context.Config.Keycloak.ClientId) {
		return nil, commonsrv.ErrPermissionDenied
	}

	return character, nil
}

func (c *characterServiceServer) validateRole(ctx context.Context, role *gocloak.Role) error {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return commonsrv.ErrPermissionDenied
	}
	if !claims.HasResourceRole(role, c.Context.Config.Keycloak.ClientId) {
		return commonsrv.ErrPermissionDenied
	}
	return nil
}

func (c *characterServiceServer) getDimension(ctx context.Context, dimensionId string) (*dimensionbus.Dimension, error) {
	dimension, err := c.Context.DimensionService.GetDimensionById(ctx, dimensionId)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrDimensionLookup, err)
		return nil, status.Error(codes.Internal, ErrDimensionLookup.Error())
	}
	if dimension == nil {
		return nil, status.Error(codes.InvalidArgument, ErrDimensionNotExist.Error())
	}
	return dimension, nil
}
