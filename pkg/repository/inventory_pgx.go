package repository

import (
	"context"
	"errors"

	"github.com/ShatteredRealms/character-service/pkg/model/inventory"
	"github.com/ShatteredRealms/go-common-service/pkg/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxInventoryRepository struct {
	conn *pgxpool.Pool
}

func NewPgxInventoryRepository(migrater *repository.PgxMigrater) InventoryRepository {
	return &pgxInventoryRepository{
		conn: migrater.Conn,
	}
}

// GetInventory implements InventoryRepository.
func (p *pgxInventoryRepository) GetInventory(ctx context.Context, characterId *uuid.UUID) (inventory.Inventory, error) {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `SELECT 
			item_id,
			slot,
			quantity
		FROM character_inventory_items
		WHERE character_id = $1;`, characterId)
	if err != nil {
		return nil, err
	}

	var inv inventory.Inventory
	inv, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[inventory.ItemInstance])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return inv, tx.Commit(ctx)
}

// SetInventory implements InventoryRepository.
func (p *pgxInventoryRepository) SetInventory(ctx context.Context, characterId *uuid.UUID, inventory *inventory.Inventory) error {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		DELETE FROM character_inventory_items 
		WHERE character_id = $1;`,
		characterId,
	)
	if err != nil {
		return err
	}
	if len(*inventory) > 0 {
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"character_inventory_items"},
			[]string{"item_id", "slot", "quantity"},
			pgx.CopyFromRows(*inventory.ToRows()),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// AddItem implements InventoryRepository.
func (p *pgxInventoryRepository) AddItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO character_inventory_items (character_id, item_id, slot, quantity)
		VALUES ($1, $2, $3, $4);`,
		characterId,
		item.Id,
		item.Slot,
		item.Quantity,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RemoveItem implements InventoryRepository.
func (p *pgxInventoryRepository) RemoveItem(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM character_inventory_items 
		WHERE character_id = $1 AND item_id = $2 AND slot = $3;`,
		characterId,
		item.Id,
		item.Slot,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// ReplaceItem implements InventoryRepository.
func (p *pgxInventoryRepository) SetQuantity(ctx context.Context, characterId *uuid.UUID, item *inventory.ItemInstance) error {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE character_inventory_items 
		SET quantity = $1
		WHERE character_id = $2 AND item_id = $3 AND slot = $4;`,
		item.Quantity,
		characterId,
		item.Id,
		item.Slot,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
