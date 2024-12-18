package service

import (
	"fmt"

	"github.com/ShatteredRealms/character-service/pkg/common"
)

var (
	ErrInvalidOwnerId = fmt.Errorf("%w: invalid owner id", common.ErrRequestInvalid)
)
