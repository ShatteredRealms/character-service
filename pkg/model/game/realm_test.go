package game_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/model/game"
)

var _ = Describe("pkg/model/game.Realm", func() {
	Describe("IsValidRealm", func() {
		It("should pass for valid realms", func() {
			Expect(game.IsValidRealm(game.RealmHuman)).To(BeTrue())
			Expect(game.IsValidRealm(game.RealmCyborg)).To(BeTrue())
		})

		It("should fail for invalid realms", func() {
			Expect(game.IsValidRealm("invalid")).To(BeFalse())
			Expect(game.IsValidRealm(game.Realm(faker.Username()))).To(BeFalse())
		})
	})
})

