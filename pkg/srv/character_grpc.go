package srv

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/gender"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/profession"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/realm"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	commonpb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/google/uuid"
	fieldmask_util "github.com/mennanov/fieldmask-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	s := &characterServiceServer{
		Context: srvCtx,
	}
	return s, nil
}

// AddCharacterPlayTime implements pb.CharacterServiceServer.
func (s *characterServiceServer) AddCharacterPlayTime(ctx context.Context, request *pb.AddPlayTimeRequest) (*emptypb.Empty, error) {
	character, err := s.Context.getCharacterById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	character, err = s.Context.CharacterService.AddCharacterPlaytime(ctx, character, request.Time)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterPlaytime, err)
		return nil, status.Error(codes.Internal, ErrCharacterPlaytime.Error())
	}

	return &emptypb.Empty{}, nil
}

// CreateCharacter implements pb.CharacterServiceServer.
func (s *characterServiceServer) CreateCharacter(ctx context.Context, request *pb.CreateCharacterRequest) (*pb.Character, error) {
	// Validate dimension exists
	dimension, err := s.Context.getDimension(ctx, request.GetDimensionId())
	if err != nil {
		return nil, err
	}

	character, err := s.Context.CharacterService.CreateCharacter(ctx, request.OwnerId, request.Name, request.Gender, request.Realm, request.Profession, &dimension.Id)
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
	character, err := s.Context.getCharacterById(ctx, request.Id)
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
func (s *characterServiceServer) EditCharacter(ctx context.Context, request *pb.EditCharacterRequest) (*pb.Character, error) {
	char, err := s.Context.getCharacterById(ctx, request.Character.Id)
	if err != nil {
		return nil, err
	}

	changeMap := make(map[string]any)
	mask, err := fieldmask_util.MaskFromProtoFieldMask(request.Mask, util.PascalCase)
	for path := range mask {
		changeMap[path] = struct{}{}
	}

	publishChanges := false
	if val, ok := changeMap["OwnerId"]; ok {
		ownerId, err := uuid.Parse(val.(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, ErrCharacterOwnerIdInvalid.Error())
		}
		char.OwnerId = ownerId
		publishChanges = true
	}
	if _, ok := changeMap["Name"]; ok {
		char.Name = request.Character.Name
	}
	if val, ok := changeMap["Gender"]; ok {
		char.Gender = val.(gender.Gender)
	}
	if val, ok := changeMap["Realm"]; ok {
		char.Realm = val.(realm.Realm)
	}
	if val, ok := changeMap["Profession"]; ok {
		char.Profession = val.(profession.Profession)
	}
	if val, ok := changeMap["PlayTime"]; ok {
		char.PlayTime = val.(int32)
	}
	if _, ok := changeMap["Location"]; ok {
		out, err := commongame.LocationFromPb(request.Character.Location)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, ErrCharcaterWorldIdInvalid.Error())
		}
		char.Location = *out
		publishChanges = true
	}
	if val, ok := changeMap["DimensionId"]; ok {
		dimension, err := s.Context.getDimension(ctx, val.(string))
		if err != nil {
			return nil, err
		}
		char.DimensionId = dimension.Id
		publishChanges = true
	}

	editedCharacter, err := s.Context.CharacterService.EditCharacter(ctx, char)
	if err != nil {
		if errors.Is(err, character.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterEdit, err)
		return nil, status.Error(codes.Internal, ErrCharacterEdit.Error())
	}

	if publishChanges {
		s.Context.CharacterBusWriter.Publish(ctx, characterbus.Message{
			Id:          editedCharacter.Id,
			OwnerId:     editedCharacter.OwnerId,
			DimensionId: editedCharacter.DimensionId,
			MapId:       editedCharacter.Location.WorldId,
			Deleted:     false,
		})
	}

	return editedCharacter.ToPb(), nil
}

// GetCharacter implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharacter(ctx context.Context, request *pb.GetCharacterRequest) (*pb.Character, error) {
	paths := []string{}
	if request.Mask != nil {
		paths = request.Mask.Paths
	}

	character, err := s.Context.getCharacterById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return character.ToPbWithMask(paths)
}

// GetCharacters implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharacters(ctx context.Context, request *pb.GetCharactersRequest) (*pb.Characters, error) {
	paths := []string{}
	if request.Mask != nil {
		paths = request.Mask.Paths
	}

	characters, _, err := s.Context.CharacterService.GetCharacters(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %e", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}

	return characters.ToPbWithMask(paths)
}

// GetCharactersForUser implements pb.CharacterServiceServer.
func (s *characterServiceServer) GetCharactersForUser(ctx context.Context, request *pb.GetUserCharactersRequest) (*pb.Characters, error) {
	characters, _, err := s.Context.CharacterService.GetCharactersByOwner(ctx, request.OwnerId)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrCharacterGet, err)
		return nil, status.Error(codes.Internal, ErrCharacterGet.Error())
	}

	return characters.ToPb(), nil
}
