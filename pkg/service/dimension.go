package service

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/repository"
)

type DimensionService interface {
	GetDimensions(ctx context.Context) (*game.Dimensions, error)
	GetDimensionById(ctx context.Context, dimensionId string) (*game.Dimension, error)
	CreateDimension(ctx context.Context, dimensionId string) (*game.Dimension, error)
	DeleteDimension(ctx context.Context, dimensionId string) (*game.Dimension, error)
}

type dimensionService struct {
	repo repository.DimensionRepository
}

func NewDimensionService(repo repository.DimensionRepository) DimensionService {
	return &dimensionService{repo: repo}
}

// CreateDimension implements DimensionService.
func (d *dimensionService) CreateDimension(ctx context.Context, dimensionId string) (*game.Dimension, error) {
	return d.repo.CreateDimension(ctx, dimensionId)
}

// DeleteDimension implements DimensionService.
func (d *dimensionService) DeleteDimension(ctx context.Context, dimensionId string) (*game.Dimension, error) {
	return d.repo.DeleteDimension(ctx, dimensionId)
}

// GetDimensionById implements DimensionService.
func (d *dimensionService) GetDimensionById(ctx context.Context, dimensionId string) (*game.Dimension, error) {
	return d.repo.GetDimensionById(ctx, dimensionId)
}

// GetDimensions implements DimensionService.
func (d *dimensionService) GetDimensions(ctx context.Context) (*game.Dimensions, error) {
	return d.repo.GetDimensions(ctx)
}
