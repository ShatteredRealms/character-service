package repository

import (
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/common"
)

var (
	ErrInvalidDimension = fmt.Errorf("%w: invalid dimension", common.ErrRequestInvalid)
	ErrNilCharacter     = fmt.Errorf("%w: character is nil", common.ErrRequestInvalid)
	ErrNilId            = fmt.Errorf("%w: id is nil", common.ErrRequestInvalid)
	ErrNonEmptyId       = fmt.Errorf("%w: id is not empty", common.ErrRequestInvalid)
)
