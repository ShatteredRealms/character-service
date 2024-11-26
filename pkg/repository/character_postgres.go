package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
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

func NewPostgresCharacterRepository(db *gorm.DB) CharacterRepository {
	db.AutoMigrate(&character.Character{})
	return &postgresCharacterRepository{gormdb: db}
}

// DeleteCharacter implements CharacterRepository.
func (p *postgresCharacterRepository) DeleteCharacter(ctx context.Context, characterId string) (character *character.Character, err error) {
	err = p.db(ctx).Clauses(clause.Returning{}).Delete(&character, "id = ?", characterId).Error
	updateSpanWithCharacter(ctx, character)
	return character, err
}

// DeleteCharactersByOwner implements CharacterRepository.
func (p *postgresCharacterRepository) DeleteCharactersByOwner(ctx context.Context, ownerId string) (characters *character.Characters, err error) {
	err = p.db(ctx).Clauses(clause.Returning{}).Delete(&characters, "owner_id = ?", ownerId).Error
	updateSpanWithOwner(ctx, ownerId)
	return characters, err
}

// GetCharacters implements CharacterRepository.
func (p *postgresCharacterRepository) GetCharacters(ctx context.Context) (characters *character.Characters, _ error) {
	return characters, p.db(ctx).Find(&characters).Error
}

// GetCharacterById implements CharacterRepository.
func (p *postgresCharacterRepository) GetCharacterById(ctx context.Context, characterId string) (*character.Character, error) {
	var character *character.Character
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
func (p *postgresCharacterRepository) GetCharactersByOwner(ctx context.Context, ownerId string) (characters *character.Characters, err error) {
	err = p.db(ctx).Where("owner_id = ?", ownerId).Find(&characters).Error
	if err != nil {
		return nil, err
	}

	updateSpanWithOwner(ctx, ownerId)
	return characters, nil
}

func (p *postgresCharacterRepository) UpdateCharacter(ctx context.Context, updatedCharacter *character.Character) (*character.Character, error) {
	if err := p.db(ctx).Save(&updatedCharacter).Error; err != nil {
		return nil, err
	}
	updateSpanWithCharacter(ctx, updatedCharacter)
	return updatedCharacter, nil
}

// CreateCharacter implements CharacterRepository.
func (p *postgresCharacterRepository) CreateCharacter(ctx context.Context, newCharacter *character.Character) (*character.Character, error) {
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

func updateSpanWithCharacter(ctx context.Context, character *character.Character) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(character.Id.String()),
		srospan.TargetCharacterName(character.Name),
	)
}
