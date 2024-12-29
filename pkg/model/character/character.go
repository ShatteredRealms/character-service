package character

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
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
	// model.Model `mapstructure:",squash"`
	Id        uuid.UUID  `db:"id" json:"id" mapstructure:"id"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" mapstructure:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt" mapstructure:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt" mapstructure:"deleted_at"`

	// Owner The account that owns the character
	OwnerId     uuid.UUID   `db:"owner_id" json:"ownerId" mapstructure:"owner_id"`
	Name        string      `db:"name" json:"name" mapstructure:"name"`
	Gender      game.Gender `db:"gender" json:"gender" mapstructure:"gender"`
	Realm       game.Realm  `db:"realm" json:"realm" mapstructure:"realm"`
	DimensionId uuid.UUID   `db:"dimension_id" json:"dimensionId" mapstructure:"dimension_id"`

	// PlayTime Time in seconds the character has played
	PlayTime int32 `db:"play_time" json:"playTime" mapstructure:"play_time"`

	// Location last location recorded for the character
	commongame.Location `json:"location" mapstructure:",squash"`
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

func (c *Character) ToPb() *pb.Character {
	char := &pb.Character{
		Id:          c.Id.String(),
		OwnerId:     c.OwnerId.String(),
		Name:        c.Name,
		Gender:      string(c.Gender),
		Realm:       string(c.Realm),
		PlayTime:    c.PlayTime,
		Location:    c.Location.ToPb(),
		DimensionId: c.DimensionId.String(),
		CreatedAt:   c.CreatedAt.Unix(),
		UpdatedAt:   c.UpdatedAt.Unix(),
	}
	if c.DeletedAt != nil {
		char.DeletedAt = c.DeletedAt.Unix()
	}
	return char
}

func (c Characters) ToPb() *pb.Characters {
	resp := &pb.Characters{Characters: make([]*pb.Character, len(c))}
	for idx, character := range c {
		if character != nil {
			resp.Characters[idx] = character.ToPb()
		}
	}

	return resp
}
