package character

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/model"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	goaway "github.com/TwiN/go-away"
)

const (
	MinCharacterNameLength = 3
	MaxCharacterNameLength = 20
)

var (
	NameRegex, _ = regexp.Compile("^[a-zA-Z0-9]+$")

	// ErrValidation thrown when a character fails validation
	ErrValidation = errors.New("validation")

	// ErrNameToShort thrown when a character name is too short
	ErrNameToShort = fmt.Errorf("%w: name must be at least %d characters", ErrValidation, MinCharacterNameLength)

	// ErrNameToLong thrown when a character name is too long
	ErrNameToLong = fmt.Errorf("%w: name can be at most %d characters", ErrValidation, MaxCharacterNameLength)

	// ErrCharacterNameToLong thrown when a character name contains invalid characters
	ErrNameInvalid = fmt.Errorf("%w: name contains invalid characters", ErrValidation)

	// ErrCharacterNameToLong thrown when a character name contains profane words
	ErrNameProfane = fmt.Errorf("%w: character name contains invalid words", ErrValidation)
)

type Character struct {
	model.Model

	// Owner The username/account that owns the character
	OwnerId   string          `gorm:"not null" json:"owner"`
	Name      string          `gorm:"not null;uniqueIndex:udx_name" json:"name"`
	Gender    game.Gender     `gorm:"not null" json:"gender"`
	Realm     game.Realm      `gorm:"not null" json:"realm"`
	Dimension *game.Dimension `gorm:"not null;uniqueIndex:udx_name;foreignKey:Id;references:Id" json:"dimension"`

	// PlayTime Time in minutes the character has played
	PlayTime uint64 `gorm:"not null" json:"play_time"`

	// Location last location recorded for the character
	Location *commongame.Location `gorm:"type:bytes;serializer:gob" json:"location"`
}
type Characters []*Character

func (c *Character) Validate() error {
	if err := c.ValidateGender(); err != nil {
		return err
	}

	if err := c.ValidateRealm(); err != nil {
		return err
	}

	return c.ValidateName()
}

func (c *Character) ValidateName() error {
	if len(c.Name) < MinCharacterNameLength {
		return ErrNameToShort
	}

	if len(c.Name) > MaxCharacterNameLength {
		return ErrNameToLong
	}

	if !NameRegex.MatchString(c.Name) {
		return ErrNameInvalid
	}

	if goaway.IsProfane(c.Name) {
		return ErrNameProfane
	}

	return nil
}

func (c *Character) ValidateGender() error {
	if game.IsValidGender(c.Gender) {
		return nil
	}

	return game.ErrorInvalidGender
}

func (c *Character) ValidateRealm() error {
	if game.IsValidRealm(c.Realm) {
		return nil
	}

	return game.ErrorInvalidRealm
}

func (c *Character) ToPb() *pb.CharacterDetails {
	return &pb.CharacterDetails{
		CharacterId: c.Id.String(),
		OwnerId:     c.OwnerId,
		Name:        c.Name,
		Gender:      string(c.Gender),
		Realm:       string(c.Realm),
		PlayTime:    c.PlayTime,
		Location:    c.Location.ToPb(),
		DimensionId: c.Dimension.Id,
	}
}

func (c Characters) ToPb() *pb.CharactersDetails {
	resp := &pb.CharactersDetails{Characters: make([]*pb.CharacterDetails, len(c))}
	for idx, character := range c {
		resp.Characters[idx] = character.ToPb()
	}

	return resp
}
