package srv

import (
	"context"
	"errors"

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

	RoleCharacterManagement = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("manage"),
		Description: gocloak.StringP("Allows creating, reading and deleting of own characters"),
	}, &CharacterRoles)

	RoleCharacterManagementOther = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("manage_other"),
		Description: gocloak.StringP("Allows creating, reading, editing and deleting of any characters"),
	}, &CharacterRoles)

	RoleAddPlaytime = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("playtime"),
		Description: gocloak.StringP("Allows adding playtime to any character"),
	}, &CharacterRoles)
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
	s := &characterServiceServer{
		Context: srvCtx,
	}
	return s, nil
}

// AddCharacterPlayTime implements pb.CharacterServiceServer.
func (s *characterServiceServer) AddCharacterPlayTime(ctx context.Context, request *pb.AddPlayTimeRequest) (*pb.PlayTimeResponse, error) {
	err := s.validateRole(ctx, RoleAddPlaytime)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetCharacterId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrCharacterIdInvalid.Error())
	}

	character, err := s.getCharacterById(ctx, &id)
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
	err := s.validateManagementPermission(ctx, request.OwnerId)
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
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrCharacterIdInvalid.Error())
	}

	_, err = s.getCharacterAndAuthCheck(ctx, &id)
	if err != nil {
		return nil, err
	}

	c, err := s.Context.CharacterService.DeleteCharacter(ctx, &id)
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
	err := s.validateRole(ctx, RoleCharacterManagementOther)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetCharacterId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrCharacterIdInvalid.Error())
	}

	char, err := s.getCharacterById(ctx, &id)
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
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrCharacterIdInvalid.Error())
	}

	character, err := s.getCharacterAndAuthCheck(ctx, &id)
	if err != nil {
		return nil, err
	}

	return character.ToPb(), nil
}

// GetCharacters implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharacters(ctx context.Context, request *emptypb.Empty) (*pb.CharactersDetails, error) {
	err := s.validateRole(ctx, RoleCharacterManagementOther)
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
	err := s.validateManagementPermission(ctx, request.Id)
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
