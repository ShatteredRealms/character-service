package repository

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/go-common-service/pkg/srospan"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresDimensionRepository struct {
	gormdb *gorm.DB
}

func NewPostgresDimensionRepository(db *gorm.DB) DimensionRepository {
	db.AutoMigrate(&game.Dimension{})
	return &postgresDimensionRepository{gormdb: db}
}

// CreateDimension implements DimensionRepository.
func (p *postgresDimensionRepository) CreateDimension(ctx context.Context, dimensionId string) (dimension *game.Dimension, _ error) {
	updateSpanWithDimension(ctx, dimensionId)
	dimension.Id = dimensionId
	return dimension, p.db(ctx).Create(dimension).Error
}

// DeleteDimension implements DimensionRepository.
func (p *postgresDimensionRepository) DeleteDimension(ctx context.Context, dimensionId string) (dimension *game.Dimension, _ error) {
	updateSpanWithDimension(ctx, dimensionId)
	return dimension, p.db(ctx).Clauses(clause.Returning{}).Delete(dimension, "id = ?", dimensionId).Error
}

// GetDimensionById implements DimensionRepository.
func (p *postgresDimensionRepository) GetDimensionById(ctx context.Context, dimensionId string) (dimension *game.Dimension, _ error) {
	result := p.db(ctx).Where("id = ?", dimensionId).Find(&dimension)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	updateSpanWithDimension(ctx, dimensionId)
	return dimension, nil
}

// GetDimensions implements DimensionRepository.
func (p *postgresDimensionRepository) GetDimensions(ctx context.Context) (dimensions *game.Dimensions, _ error) {
	return dimensions, p.db(ctx).Find(dimensions).Error
}

func (p *postgresDimensionRepository) db(ctx context.Context) *gorm.DB {
	return p.gormdb.WithContext(ctx)
}

func updateSpanWithDimension(ctx context.Context, dimensionId string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.DimensionId(dimensionId),
	)
}
