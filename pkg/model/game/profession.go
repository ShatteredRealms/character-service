package game

import (
	"errors"
	"strings"
)

type Profession string

const (
	ProfessionNecromancer Profession = "Necromancer"
)

var (
	ErrorInvalidProfession = errors.New("invalid profession")
)

func FromString(s string) (Profession, error) {
	s = strings.ToLower(s)
	switch s {
	case "necromancer":
		return ProfessionNecromancer, nil
	default:
		return "", ErrorInvalidProfession
	}
}
