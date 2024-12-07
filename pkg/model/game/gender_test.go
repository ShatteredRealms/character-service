package game_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
)

var _ = Describe("pkg/model/game.Gender", func() {
	Describe("IsValidGender", func() {
		It("should pass for valid genders", func() {
			Expect(game.IsValidGender(game.GenderMale)).To(BeTrue())
			Expect(game.IsValidGender(game.GenderFemale)).To(BeTrue())
		})

		It("should fail for invalid genders", func() {
			Expect(game.IsValidGender("invalid")).To(BeFalse())
			Expect(game.IsValidGender(game.Gender(faker.Username()))).To(BeFalse())
		})
	})
})
