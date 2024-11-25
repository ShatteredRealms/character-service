package repository

import (
	"context"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
)

type DimensionRepository interface {
	GetDimensionById(ctx context.Context, dimensionId string) (*game.Dimension, error)

	GetDimensions(ctx context.Context) (*game.Dimensions, error)

	CreateDimension(ctx context.Context, dimensionId string) (*game.Dimension, error)

	DeleteDimension(ctx context.Context, dimensionId string) (*game.Dimension, error)
}
