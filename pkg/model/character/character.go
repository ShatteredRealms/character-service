package character

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/gender"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/profession"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/realm"
	"github.com/ShatteredRealms/go-common-service/pkg/model/game"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	goaway "github.com/TwiN/go-away"
	"github.com/google/uuid"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
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
	OwnerId     uuid.UUID             `db:"owner_id" json:"ownerId" mapstructure:"owner_id"`
	Name        string                `db:"name" json:"name" mapstructure:"name"`
	Gender      gender.Gender         `db:"gender" json:"gender" mapstructure:"gender"`
	Realm       realm.Realm           `db:"realm" json:"realm" mapstructure:"realm"`
	Profession  profession.Profession `db:"profession" json:"profession" mapstructure:"profession"`
	DimensionId uuid.UUID             `db:"dimension_id" json:"dimensionId" mapstructure:"dimension_id"`

	// PlayTime Time in seconds the character has played
	PlayTime int32 `db:"play_time" json:"playTime" mapstructure:"play_time"`

	// SkillStats map[skills.SkillName]int

	// Location last location recorded for the character
	game.Location `json:"location" mapstructure:",squash"`
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

	if err := c.ValidateProfession(); err != nil {
		return err
	}

	return c.ValidateName()
}

func (c *Character) ValidateProfession() error {
	_, err := profession.FromString(string(c.Profession))
	return err
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
	if gender.IsValid(c.Gender) {
		return nil
	}

	return gender.ErrorInvalid
}

func (c *Character) ValidateRealm() error {
	if realm.IsValid(c.Realm) {
		return nil
	}

	return realm.ErrorInvalidRealm
}

func (c *Character) ToPb() *pb.Character {
	char := &pb.Character{
		Id:          c.Id.String(),
		OwnerId:     c.OwnerId.String(),
		Name:        c.Name,
		Gender:      string(c.Gender),
		Realm:       string(c.Realm),
		Profession:  string(c.Profession),
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

func (c *Character) ToPbWithMask(paths []string) (*pb.Character, error) {
	mask, err := fieldmask_utils.MaskFromPaths(paths, util.PascalCase)
	if err != nil {
		return nil, err
	}

	outPb := &pb.Character{}
	err = fieldmask_utils.StructToStruct(mask, c.ToPb(), outPb)
	return outPb, err
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

func (c Characters) ToPbWithMask(paths []string) (*pb.Characters, error) {
	var err error
	resp := &pb.Characters{Characters: make([]*pb.Character, len(c))}
	for idx, character := range c {
		if character != nil {
			resp.Characters[idx], err = character.ToPbWithMask(paths)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}
