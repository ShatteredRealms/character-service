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

type Model struct {
	model.Model

	// Owner The username/account that owns the character
	OwnerId     string                  `gorm:"not null" json:"owner"`
	Name        string                  `gorm:"uniqueIndex:idx_deleted" json:"name"`
	Gender      game.Gender             `gorm:"not null" json:"gender"`
	Realm       game.Realm              `gorm:"not null" json:"realm"`
	DimensionId string                  `gorm:"not null" json:"dimensionId"`
	Dimension   *dimensionbus.Dimension `gorm:"not null" json:"dimension"`

	// PlayTime Time in minutes the character has played
	PlayTime uint64 `gorm:"not null" json:"play_time"`

	// Location last location recorded for the character
	Location commongame.Location `gorm:"type:bytes;serializer:gob" json:"location"`
}
type Models []*Model

var (
	detector *goaway.ProfanityDetector
)

func init() {
	detector = goaway.NewProfanityDetector()
	detector.WithSanitizeLeetSpeak(true)
}

func (c *Model) Validate() error {
	if err := c.ValidateGender(); err != nil {
		return err
	}

	if err := c.ValidateRealm(); err != nil {
		return err
	}

	return c.ValidateName()
}

func (c *Model) ValidateName() error {
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

func (c *Model) ValidateGender() error {
	if game.IsValidGender(c.Gender) {
		return nil
	}

	return game.ErrorInvalidGender
}

func (c *Model) ValidateRealm() error {
	if game.IsValidRealm(c.Realm) {
		return nil
	}

	return game.ErrorInvalidRealm
}

func (c *Model) ToPb() *pb.CharacterDetails {
	return &pb.CharacterDetails{
		CharacterId: c.Id.String(),
		OwnerId:     c.OwnerId,
		Name:        c.Name,
		Gender:      string(c.Gender),
		Realm:       string(c.Realm),
		PlayTime:    c.PlayTime,
		Location:    c.Location.ToPb(),
		DimensionId: c.DimensionId,
	}
}

func (c Models) ToPb() *pb.CharactersDetails {
	resp := &pb.CharactersDetails{Characters: make([]*pb.CharacterDetails, len(c))}
	for idx, character := range c {
		if character != nil {
			resp.Characters[idx] = character.ToPb()
		}
	}

	return resp
}
