package srv

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/WilSimpson/gocloak/v13"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	CharacterRoles = make([]*gocloak.Role, 0)

	RoleCharacterManagement = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("manage"),
		Description: gocloak.StringP("Allows creating, reading and deleting of own characters"),
	}, CharacterRoles)

	RoleCharacterManagementOther = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("manage_other"),
		Description: gocloak.StringP("Allows creating, reading, editing and deleting of any characters"),
	}, CharacterRoles)
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
	return &characterServiceServer{
		Context: srvCtx,
	}, nil
}

// AddCharacterPlayTime implements pb.CharacterServiceServer.
func (c *characterServiceServer) AddCharacterPlayTime(ctx context.Context, request *pb.AddPlayTimeRequest) (*pb.PlayTimeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method AddCharacterPlayTime not implemented")
}

// CreateCharacter implements pb.CharacterServiceServer.
func (c *characterServiceServer) CreateCharacter(ctx context.Context, request *pb.CreateCharacterRequest) (*pb.CharacterDetails, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateCharacter not implemented")
}

// DeleteCharacter implements pb.CharacterServiceServer.
func (c *characterServiceServer) DeleteCharacter(ctx context.Context, request *pb.CharacterTarget) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteCharacter not implemented")
}

// EditCharacter implements pb.CharacterServiceServer.
func (c *characterServiceServer) EditCharacter(ctx context.Context, request *pb.EditCharacterRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method EditCharacter not implemented")
}

// GetCharacter implements pb.CharacterServiceServer.
func (c *characterServiceServer) GetCharacter(ctx context.Context, request *pb.CharacterTarget) (*pb.CharacterDetails, error) {
	return nil, status.Error(codes.Unimplemented, "method GetCharacter not implemented")
}

// GetCharacters implements pb.CharacterServiceServer.
func (c *characterServiceServer) GetCharacters(ctx context.Context, request *emptypb.Empty) (*pb.CharactersDetails, error) {
	return nil, status.Error(codes.Unimplemented, "method GetCharacters not implemented")
}

// GetCharactersForUser implements pb.CharacterServiceServer.
func (c *characterServiceServer) GetCharactersForUser(ctx context.Context, request *pb.UserTarget) (*pb.CharactersDetails, error) {
	return nil, status.Error(codes.Unimplemented, "method GetCharactersForUser not implemented")
}
