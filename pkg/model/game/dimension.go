package game

import (
	"errors"
	"time"
)

var (
	ErrInvalidDimension = errors.New("invalid dimension")
)

type Dimension struct {
	Id        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Dimensions []*Dimension
