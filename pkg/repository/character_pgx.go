package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ShatteredRealms/character-service/pkg/common"
	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
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

func QueryCharacters(ctx context.Context, tx pgx.Tx, matchFilters map[string]interface{}, queryFilter *pb.QueryFilters, deleted bool) (pgx.Rows, int, error) {
	total := -1
	builder := strings.Builder{}
	params := make([]interface{}, 0, len(matchFilters))
	builder.WriteString("FROM characters WHERE (deleted_at IS ")
	if deleted {
		builder.WriteString("NOT NULL")
	} else {
		builder.WriteString("NULL")
	}
	for key, value := range matchFilters {
		builder.WriteString(" AND ")
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(fmt.Sprintf("$%d", len(params)+1))
		params = append(params, value)
	}
	builder.WriteString(")")

	if queryFilter != nil {
		rows, err := tx.Query(ctx, "SELECT COUNT(*) "+builder.String(), params...)
		if err != nil {
			return nil, -1, err
		}
		if rows.Next() {
			err = rows.Scan(&total)
			rows.Close()
		}

		log.Logger.Infof("total: %d", total)

		if queryFilter.Limit > 0 {
			builder.WriteString(fmt.Sprintf(" LIMIT %d", queryFilter.Limit))
		}
		if queryFilter.Offset > 0 {
			builder.WriteString(fmt.Sprintf(" OFFSET %d", queryFilter.Offset))
		}
	}
	rows, err := tx.Query(ctx, "SELECT * "+builder.String()+";", params...)
	return rows, total, err
}

// GetCharacter implements CharacterRepository.
func (p *pgxCharacterRepository) GetCharacter(ctx context.Context, matchFilters map[string]interface{}) (*character.Character, error) {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	rows, _, err := QueryCharacters(ctx, tx, matchFilters, nil, false)
	if err != nil {
		return nil, err
	}
	outCharacter, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return outCharacter, tx.Commit(ctx)
}

// GetCharacters implements CharacterRepository.
func (p *pgxCharacterRepository) GetCharacters(ctx context.Context, matchFilters map[string]interface{}, queryFilter *pb.QueryFilters, deleted bool) (character.Characters, int, error) {
	tx, err := p.conn.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, -1, err
	}

	rows, total, err := QueryCharacters(ctx, tx, matchFilters, queryFilter, deleted)
	if err != nil {
		return nil, -1, err
	}
	outCharacters, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[character.Character])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, nil
		}
		return nil, -1, err
	}

	return outCharacters, total, tx.Commit(ctx)
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
			owner_id = $1,
			name = $2,
			gender = $3,
			realm = $4,
			dimension_id = $5,
			play_time = $6,
			world_id = $7,
			x = $8,
			y = $9,
			z = $10,
			roll = $11,
			pitch = $12,
			yaw = $13
		WHERE id = $14 RETURNING *`,
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
		c.Id,
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
