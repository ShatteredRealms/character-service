package game

import "errors"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

var (
	ErrorInvalidGender = errors.New("invalid gender")
)

func IsValidGender(gender Gender) bool {
	return gender == GenderMale ||
		gender == GenderFemale
}
