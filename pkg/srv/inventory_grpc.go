package srv

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	character, err := s.Context.getCharacterById(ctx, request.CharacterId)
	if err != nil {
		return nil, err
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
	s := &inventoryServiceServer{
		Context: srvCtx,
	}
	return s, nil
}
