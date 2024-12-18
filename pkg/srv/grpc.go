package srv

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	commonpb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	CharacterRoles = make([]*gocloak.Role, 0)

	RolePlaytime = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("character.playtime"),
		Description: gocloak.StringP("Allows modifying playtime of characters"),
	}, &CharacterRoles)

	RoleGetCharactersSelf = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("characters.get.self"),
		Description: gocloak.StringP("Allows getting a character if the user is the owner"),
	}, &CharacterRoles)

	RoleGetCharactersAll = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("character.get.all"),
		Description: gocloak.StringP("Allows getting a character even if the user is not the owner"),
	}, &CharacterRoles)

	RoleCreateCharactersSelf = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("characters.create.self"),
		Description: gocloak.StringP("Allows creating a character that the user will owner"),
	}, &CharacterRoles)

	RoleCreateCharactersAll = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("character.create.all"),
		Description: gocloak.StringP("Allows creating a character even if the owner is not the user"),
	}, &CharacterRoles)

	RoleDeleteCharactersSelf = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("characters.delete.self"),
		Description: gocloak.StringP("Allows deleting a character if the requester is the owner"),
	}, &CharacterRoles)

	RoleDeleteCharactersAll = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("character.delete.all"),
		Description: gocloak.StringP("Allows deleting a character even if the user is not the owner"),
	}, &CharacterRoles)

	RoleEditCharacter = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("character.edit"),
		Description: gocloak.StringP("Allows editing any characters details except playtime"),
	}, &CharacterRoles)
)

var (
	CompositeCharacterRoles = make([]*gocloak.Role, 0)

	RoleManageCharactersSelf = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("characters.manage.self"),
		Description: gocloak.StringP("Allows managing a character if the user is the owner"),
		Composite:   gocloak.BoolP(true),
		Composites: &gocloak.CompositesRepresentation{
			Client: &map[string][]string{
				config.GlobalConfig.Keycloak.ClientId: {"characters.get.self", "characters.create.self", "characters.delete.self"},
			},
		},
	}, &CompositeCharacterRoles)

	RoleManageCharactersAll = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("characters.manage.all"),
		Description: gocloak.StringP("Allows managing a character even if the user is not the owner and editing character details"),
		Composite:   gocloak.BoolP(true),
		Composites: &gocloak.CompositesRepresentation{
			Client: &map[string][]string{
				config.GlobalConfig.Keycloak.ClientId: {"character.get.all", "character.create.all", "character.delete.all", "character.edit"},
			},
		},
	}, &CompositeCharacterRoles)
)

var (
	ErrCharacterDoesNotExist   = errors.New("CS-C-00")
	ErrCharacterCreate         = errors.New("CS-C-01")
	ErrCharacterDelete         = errors.New("CS-C-02")
	ErrCharacterGet            = errors.New("CS-C-03")
	ErrCharacterEdit           = errors.New("CS-C-04")
	ErrCharacterPlaytime       = errors.New("CS-C-05")
	ErrCharacterIdInvalid      = errors.New("CS-C-07")
	ErrCharacterOwnerIdInvalid = errors.New("CS-C-08")
	ErrCharcaterWorldIdInvalid = errors.New("CS-C-09")

	ErrDimensionNotExist = errors.New("CS-D-01")
	ErrDimensionLookup   = errors.New("CS-D-02")
)

type characterServiceServer struct {
	pb.UnimplementedCharacterServiceServer
	Context *CharacterContext
}

func NewCharacterServiceServer(ctx context.Context, srvCtx *CharacterContext) (pb.CharacterServiceServer, error) {
	err := srvCtx.CreateRoles(ctx, &CharacterRoles)
	if err != nil {
		return nil, err
	}
	err = srvCtx.CreateRoles(ctx, &CompositeCharacterRoles)
	if err != nil {
		return nil, err
	}

	s := &characterServiceServer{
		Context: srvCtx,
	}
	return s, nil
}

// AddCharacterPlayTime implements pb.CharacterServiceServer.
func (s *characterServiceServer) AddCharacterPlayTime(ctx context.Context, request *pb.AddPlayTimeRequest) (*pb.PlayTimeResponse, error) {
	character, err := s.validateCharacterPermissions(ctx, request.GetCharacterId(), RolePlaytime, RolePlaytime)
	if err != nil {
		return nil, err
	}

	character, err = s.Context.CharacterService.AddCharacterPlaytime(ctx, character, request.Time)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterPlaytime, err)
		return nil, status.Error(codes.Internal, ErrCharacterPlaytime.Error())
	}

	return &pb.PlayTimeResponse{
		Time: character.PlayTime,
	}, nil
}

// CreateCharacter implements pb.CharacterServiceServer.
func (s *characterServiceServer) CreateCharacter(ctx context.Context, request *pb.CreateCharacterRequest) (*pb.CharacterDetails, error) {
	err := s.validateUserPermissions(ctx, request.OwnerId, RoleCreateCharactersSelf, RoleCreateCharactersAll)
	if err != nil {
		return nil, err
	}

	// Validate dimension exists
	dimension, err := s.getDimension(ctx, request.GetDimensionId())
	if err != nil {
		return nil, err
	}

	character, err := s.Context.CharacterService.CreateCharacter(ctx, request.OwnerId, request.Name, request.Gender, request.Realm, &dimension.Id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterCreate, err)
		return nil, status.Error(codes.Internal, ErrCharacterCreate.Error())
	}

	s.Context.CharacterBusWriter.Publish(ctx, characterbus.Message{
		Id:          character.Id,
		OwnerId:     character.OwnerId,
		DimensionId: character.DimensionId,
		MapId:       character.Location.WorldId,
		Deleted:     false,
	})

	return character.ToPb(), nil
}

// DeleteCharacter implements pb.CharacterServiceServer.
func (s *characterServiceServer) DeleteCharacter(ctx context.Context, request *commonpb.TargetId) (*emptypb.Empty, error) {
	character, err := s.validateCharacterPermissions(ctx, request.Id, RoleDeleteCharactersSelf, RoleDeleteCharactersAll)
	if err != nil {
		return nil, err
	}

	c, err := s.Context.CharacterService.DeleteCharacter(ctx, &character.Id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterDelete, err)
		return nil, status.Error(codes.Internal, ErrCharacterDelete.Error())
	}

	s.Context.CharacterBusWriter.Publish(ctx, characterbus.Message{
		Id:      c.Id,
		OwnerId: c.OwnerId,
		Deleted: true,
	})

	return &emptypb.Empty{}, nil
}

// EditCharacter implements pb.CharacterServiceServer.
func (s *characterServiceServer) EditCharacter(ctx context.Context, request *pb.EditCharacterRequest) (*emptypb.Empty, error) {
	char, err := s.validateCharacterPermissions(ctx, request.GetCharacterId(), RoleEditCharacter, RoleEditCharacter)
	if err != nil {
		return nil, err
	}

	publishChanges := false
	if request.OptionalOwnerId != nil {
		ownerId, err := uuid.Parse(request.GetOwnerId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, ErrCharacterOwnerIdInvalid.Error())
		}
		char.OwnerId = ownerId
		publishChanges = true
	}
	if request.OptionalNewName != nil {
		char.Name = request.GetNewName()
	}
	if request.OptionalGender != nil {
		char.Gender = game.Gender(request.GetGender())
	}
	if request.OptionalRealm != nil {
		char.Realm = game.Realm(request.GetRealm())
	}
	if request.OptionalPlayTime != nil {
		err := s.validateRole(ctx, RolePlaytime)
		if err != nil {
			return nil, err
		}
		char.PlayTime = request.GetPlayTime()
	}
	if request.OptionalLocation != nil {
		out, err := commongame.LocationFromPb(request.GetLocation())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, ErrCharcaterWorldIdInvalid.Error())
		}
		char.Location = *out
		publishChanges = true
	}
	if request.OptionalDimension != nil {
		dimension, err := s.getDimension(ctx, request.GetDimensionId())
		if err != nil {
			return nil, err
		}
		char.DimensionId = dimension.Id
		publishChanges = true
	}

	c, err := s.Context.CharacterService.EditCharacter(ctx, char)
	if err != nil {
		if errors.Is(err, character.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterEdit, err)
		return nil, status.Error(codes.Internal, ErrCharacterEdit.Error())
	}

	if publishChanges {
		s.Context.CharacterBusWriter.Publish(ctx, characterbus.Message{
			Id:          c.Id,
			OwnerId:     c.OwnerId,
			DimensionId: c.DimensionId,
			MapId:       c.Location.WorldId,
			Deleted:     false,
		})
	}

	return &emptypb.Empty{}, nil
}

// GetCharacter implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharacter(ctx context.Context, request *commonpb.TargetId) (*pb.CharacterDetails, error) {
	character, err := s.validateCharacterPermissions(ctx, request.Id, RoleGetCharactersSelf, RoleGetCharactersAll)
	if err != nil {
		return nil, err
	}

	return character.ToPb(), nil
}

// GetCharacters implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharacters(ctx context.Context, request *emptypb.Empty) (*pb.CharactersDetails, error) {
	err := s.validateRole(ctx, RoleGetCharactersAll)
	if err != nil {
		return nil, err
	}

	characters, err := s.Context.CharacterService.GetCharacters(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %e", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}

	return characters.ToPb(), nil
}

// GetCharactersForUser implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharactersForUser(ctx context.Context, request *commonpb.TargetId) (*pb.CharactersDetails, error) {
	err := s.validateUserPermissions(ctx, request.Id, RoleGetCharactersSelf, RoleGetCharactersAll)
	if err != nil {
		return nil, err
	}

	characters, err := s.Context.CharacterService.GetCharactersByOwner(ctx, request.Id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}

	return characters.ToPb(), nil
}
