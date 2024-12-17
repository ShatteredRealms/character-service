package character

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	"github.com/ShatteredRealms/go-common-service/pkg/model"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	goaway "github.com/TwiN/go-away"
	"github.com/google/uuid"
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
	OwnerId     uuid.UUID               `gorm:"not null" json:"owner"`
	Name        string                  `gorm:"not null;uniqueIndex:idx_deleted" json:"name"`
	Gender      game.Gender             `gorm:"not null" json:"gender"`
	Realm       game.Realm              `gorm:"not null" json:"realm"`
	DimensionId uuid.UUID               `gorm:"not null" json:"dimensionId"`
	Dimension   *dimensionbus.Dimension `gorm:"not null" json:"dimension"`

	// PlayTime Time in minutes the character has played
	PlayTime uint64 `gorm:"not null" json:"play_time"`

	// Location last location recorded for the character
	Location commongame.Location `gorm:"type:bytes;serializer:gob" json:"location"`
}
type Characters []*Character

var (
	detector *goaway.ProfanityDetector
)

func init() {
	detector = goaway.NewProfanityDetector()
	detector.WithSanitizeLeetSpeak(true)
}

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

	if detector.IsProfane(c.Name) {
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
		OwnerId:     c.OwnerId.String(),
		Name:        c.Name,
		Gender:      string(c.Gender),
		Realm:       string(c.Realm),
		PlayTime:    c.PlayTime,
		Location:    c.Location.ToPb(),
		DimensionId: c.DimensionId.String(),
		CreatedAt:   uint64(c.CreatedAt.Unix()),
	}
}

func (c Characters) ToPb() *pb.CharactersDetails {
	resp := &pb.CharactersDetails{Characters: make([]*pb.CharacterDetails, len(c))}
	for idx, character := range c {
		if character != nil {
			resp.Characters[idx] = character.ToPb()
		}
	}

	return resp
}
