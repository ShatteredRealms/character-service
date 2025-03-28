package srv

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/WilSimpson/gocloak/v13"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InventoryRoles = make([]*gocloak.Role, 0)

	RoleGetInventorySelf = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("inventory.get.self"),
		Description: gocloak.StringP("Allows getting a character if the user is the owner"),
	}, &InventoryRoles)

	RoleGetInventoryAll = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("inventory.get.all"),
		Description: gocloak.StringP("Allows getting a character even if the user is not the owner"),
	}, &InventoryRoles)

	RoleUpdateInventory = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("inventory.edit"),
		Description: gocloak.StringP("Allows editing a character's inventory"),
	}, &InventoryRoles)
)

var (
	CompositeInventoryRoles = make([]*gocloak.Role, 0)

	RoleManageAllInventories = util.RegisterRole(&gocloak.Role{
		Name:        gocloak.StringP("inventory.all"),
		Description: gocloak.StringP("Allows managing all character inventories."),
		Composite:   gocloak.BoolP(true),
		Composites: &gocloak.CompositesRepresentation{
			Client: &map[string][]string{
				config.GlobalConfig.Keycloak.ClientId: {"inventory.get.self", "inventory.get.all", "inventory.edit"},
			},
		},
	}, &CompositeInventoryRoles)
)

var (
	ErrInventoryGet  = errors.New("IS-I-03")
	ErrInventoryEdit = errors.New("IS-I-04")
)

type inventoryServiceServer struct {
	pb.UnimplementedInventoryServiceServer
	Context *CharacterContext
}

// GetInventory implements pb.InventoryServiceServer.
func (s *inventoryServiceServer) GetInventory(
	ctx context.Context,
	request *pb.GetInventoryRequest,
) (*pb.GetInventoryResponse, error) {
	character, err := s.Context.validateCharacterPermissions(ctx, request.CharacterId, RoleGetInventorySelf, RoleGetInventoryAll)
	if err != nil {
		return nil, err
	}

	if character == nil {
		return nil, status.Error(codes.NotFound, ErrCharacterDoesNotExist.Error())
	}

	inventory, err := s.Context.InventoryService.GetInventory(ctx, &character.Id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("code %v: %v", ErrInventoryGet, err)
		return nil, status.Error(codes.Internal, ErrInventoryGet.Error())
	}

	return &pb.GetInventoryResponse{
		Items: inventory.ToPb(),
	}, nil
}

func NewInventoryServiceServer(ctx context.Context, srvCtx *CharacterContext) (pb.InventoryServiceServer, error) {
	err := srvCtx.CreateRoles(ctx, &InventoryRoles)
	if err != nil {
		return nil, err
	}
	err = srvCtx.CreateRoles(ctx, &CompositeInventoryRoles)
	if err != nil {
		return nil, err
	}

	s := &inventoryServiceServer{
		Context: srvCtx,
	}
	return s, nil
}
