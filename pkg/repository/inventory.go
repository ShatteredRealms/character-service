package repository

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/inventory"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	// GetCharacters gets the of character inventory for the provided character id
	GetInventory(ctx context.Context, characterId *uuid.UUID) (inventory.Inventory, error)
	SetInventory(ctx context.Context, characterId *uuid.UUID, inventory *inventory.Inventory) error
	AddItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error
	RemoveItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error
	SetQuantity(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error
}
