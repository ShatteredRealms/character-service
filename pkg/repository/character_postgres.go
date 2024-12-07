package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrCharacter           = errors.New("character repository")
	ErrCharacterIdProvided = fmt.Errorf("%w: id provided", ErrCharacter)
)

type postgresCharacterRepository struct {
	gormdb *gorm.DB
}

func NewPostgresCharacter(db *gorm.DB) (CharacterRepository, error) {
	return &postgresCharacterRepository{gormdb: db}, db.AutoMigrate(&character.Model{})
}

// DeleteCharacter implements CharacterRepository.
func (p *postgresCharacterRepository) DeleteCharacter(ctx context.Context, characterId *uuid.UUID) (character *character.Model, _ error) {
	result := p.db(ctx).Clauses(clause.Returning{}).Delete(&character, "id = ?", characterId)
	if result.RowsAffected > 0 {
		updateSpanWithCharacter(ctx, character)
	} else {
		character = nil
	}
	return character, result.Error
}

// DeleteCharactersByOwner implements CharacterRepository.
func (p *postgresCharacterRepository) DeleteCharactersByOwner(ctx context.Context, ownerId string) (characters *character.Models, err error) {
	err = p.db(ctx).Clauses(clause.Returning{}).Delete(&characters, "owner_id = ?", ownerId).Error
	updateSpanWithOwner(ctx, ownerId)
	return characters, err
}

// GetCharacters implements CharacterRepository.
func (p *postgresCharacterRepository) GetCharacters(ctx context.Context) (characters *character.Models, _ error) {
	return characters, p.db(ctx).Find(&characters).Error
}

// GetCharacterById implements CharacterRepository.
func (p *postgresCharacterRepository) GetCharacterById(ctx context.Context, characterId *uuid.UUID) (*character.Model, error) {
	var character *character.Model
	result := p.db(ctx).Where("id = ?", characterId).Find(&character)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	updateSpanWithCharacter(ctx, character)
	return character, nil
}

// GetCharactersByOwner implements CharacterRepository.
func (p *postgresCharacterRepository) GetCharactersByOwner(ctx context.Context, ownerId string) (characters *character.Models, err error) {
	updateSpanWithOwner(ctx, ownerId)
	return characters, p.db(ctx).Where("owner_id = ?", ownerId).Find(&characters).Error
}

func (p *postgresCharacterRepository) UpdateCharacter(ctx context.Context, updatedCharacter *character.Model) (*character.Model, error) {
	if err := p.db(ctx).Save(&updatedCharacter).Error; err != nil {
		return nil, err
	}
	updateSpanWithCharacter(ctx, updatedCharacter)
	return updatedCharacter, nil
}

// CreateCharacter implements CharacterRepository.
func (p *postgresCharacterRepository) CreateCharacter(ctx context.Context, newCharacter *character.Model) (*character.Model, error) {
	if newCharacter.Id != nil {
		return nil, ErrCharacterIdProvided
	}

	if err := p.db(ctx).Create(&newCharacter).Error; err != nil {
		return nil, err
	}

	updateSpanWithCharacter(ctx, newCharacter)
	return newCharacter, nil
}

func (p *postgresCharacterRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithOwner(ctx context.Context, ownerId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetOwnerId(ownerId),
	)
}

func updateSpanWithCharacter(ctx context.Context, character *character.Model) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(character.Id.String()),
		srospan.TargetCharacterName(character.Name),
	)
}
