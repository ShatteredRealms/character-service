package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/common"
	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxCharacterRepository struct {
	conn *pgxpool.Pool
}

func NewPgxCharacterRepository(migrater *repository.PgxMigrater) CharacterRepository {
	return &pgxCharacterRepository{
		conn: migrater.Conn,
	}
}

// CreateCharacter implements CharacterRepository.
func (p *pgxCharacterRepository) CreateCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error) {
	if newCharacter == nil {
		return nil, ErrNilCharacter
	}

	if (newCharacter.Id != uuid.UUID{}) {
		return nil, ErrNonEmptyId
	}

	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(
		ctx,
		`INSERT INTO characters (
			owner_id,
			dimension_id,
			name,
			gender,
			realm,
			world_id,
			x,
			y,
			z,
			roll,
			pitch,
			yaw
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *`,
		newCharacter.OwnerId,
		newCharacter.DimensionId,
		newCharacter.Name,
		newCharacter.Gender,
		newCharacter.Realm,
		newCharacter.Location.WorldId,
		newCharacter.Location.X,
		newCharacter.Location.Y,
		newCharacter.Location.Z,
		newCharacter.Location.Roll,
		newCharacter.Location.Pitch,
		newCharacter.Location.Yaw,
	)
	outCharacter, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		return nil, err
	}

	return outCharacter, tx.Commit(ctx)
}

// DeleteCharacter implements CharacterRepository.
func (p *pgxCharacterRepository) DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	if characterId == nil {
		return nil, ErrNilId
	}

	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx,
		"UPDATE characters SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING *",
		characterId)
	outCharacter, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacter, tx.Commit(ctx)
}

// DeleteCharactersByOwner implements CharacterRepository.
func (p *pgxCharacterRepository) DeleteCharactersByOwner(ctx context.Context, ownerId *uuid.UUID) (character.Characters, error) {
	if ownerId == nil {
		return nil, ErrNilId
	}

	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx,
		"UPDATE characters SET deleted_at	= CURRENT_TIMESTAMP WHERE owner_id = $1 RETURNING *",
		ownerId)
	outCharacters, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacters, tx.Commit(ctx)
}

// GetCharacterById implements CharacterRepository.
func (p *pgxCharacterRepository) GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Character, error) {
	if characterId == nil {
		return nil, ErrNilId
	}

	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx,
		"SELECT * FROM characters WHERE id = $1",
		characterId)
	outCharacter, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacter, tx.Commit(ctx)
}

// GetCharacters implements CharacterRepository.
func (p *pgxCharacterRepository) GetCharacters(ctx context.Context) (character.Characters, error) {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx, "SELECT * FROM characters WHERE deleted_at IS NULL")
	outCharacters, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacters, tx.Commit(ctx)
}

// GetCharactersByOwner implements CharacterRepository.
func (p *pgxCharacterRepository) GetCharactersByOwner(ctx context.Context, ownerId *uuid.UUID) (character.Characters, error) {
	if ownerId == nil {
		return nil, ErrNilId
	}

	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx,
		"SELECT * FROM characters WHERE (owner_id = $1 AND deleted_at IS NULL)",
		ownerId)
	outCharacters, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacters, tx.Commit(ctx)
}

// GetDeletedCharacters implements CharacterRepository.
func (p *pgxCharacterRepository) GetDeletedCharacters(ctx context.Context) (character.Characters, error) {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx, "SELECT * FROM characters WHERE deleted_at IS NOT NULL")
	outCharacters, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacters, tx.Commit(ctx)
}

// UpdateCharacter implements CharacterRepository.
func (p *pgxCharacterRepository) UpdateCharacter(ctx context.Context, c *character.Character) (*character.Character, error) {
	if c == nil {
		return nil, ErrNilCharacter
	}

	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := tx.Query(ctx,
		`UPDATE 
			characters 
		SET 
			owner_id = $2,
			name = $3,
			gender = $4,
			realm = $5,
			dimension_id = $6,
			play_time = $7,
			world_id = $8,
			x = $9,
			y = $10,
			z = $11,
			roll = $12,
			pitch = $13,
			yaw = $14
		WHERE id = $1 RETURNING *`,
		c.Id,
		c.OwnerId,
		c.Name,
		c.Gender,
		c.Realm,
		c.DimensionId,
		c.PlayTime,
		c.Location.WorldId,
		c.Location.X,
		c.Location.Y,
		c.Location.Z,
		c.Location.Roll,
		c.Location.Pitch,
		c.Location.Yaw,
	)
	outCharacter, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", common.ErrRequestInvalid, err)
		}
		return nil, err
	}

	return outCharacter, tx.Commit(ctx)
}
