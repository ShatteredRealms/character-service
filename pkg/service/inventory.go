package service

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/model/inventory"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	"github.com/google/uuid"
)

var (
	ErrInventory = errors.New("inventory service")
)

type InventoryService interface {
	GetInventory(ctx context.Context, characterId *uuid.UUID) (inventory.Inventory, error)
	SetInventory(ctx context.Context, characterId *uuid.UUID, inventory *inventory.Inventory) error
	AddItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error
	RemoveItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error
	SetQuantity(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

// AddItem implements InventoryService.
func (i *inventoryService) AddItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error {
	return i.repo.AddItem(ctx, characterId, item)
}

// GetInventory implements InventoryService.
func (i *inventoryService) GetInventory(ctx context.Context, characterId *uuid.UUID) (inventory.Inventory, error) {
	return i.repo.GetInventory(ctx, characterId)
}

// RemoveItem implements InventoryService.
func (i *inventoryService) RemoveItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error {
	return i.repo.RemoveItem(ctx, characterId, item)
}

// SetInventory implements InventoryService.
func (i *inventoryService) SetInventory(ctx context.Context, characterId *uuid.UUID, inventory *inventory.Inventory) error {
	return i.repo.SetInventory(ctx, characterId, inventory)
}

// SetQuantity implements InventoryService.
func (i *inventoryService) SetQuantity(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error {
	return i.repo.SetQuantity(ctx, characterId, item)
}
