package character_test

import (
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
)

var _ = Describe("pkg/model/character.Character", func() {
	var (
		c *character.Character
	)

	BeforeEach(func() {
		c = &character.Character{
			Name:   faker.Username(),
			Gender: game.GenderMale,
			Realm:  game.RealmHuman,
		}
		id, err := uuid.NewV7()
		Expect(err).NotTo(HaveOccurred())
		c.Id = &id
		id, err = uuid.NewV7()
		Expect(err).NotTo(HaveOccurred())
		c.OwnerId = id
	})

	Describe("proto conversions", func() {
		var (
			pbs        []*pb.CharacterDetails
			characters character.Characters
		)
		It("should convert to a character to a pb ", func() {
			pbs = append(pbs, c.ToPb())
			characters = append(characters, c)
		})
		It("should convert to characters to a pb", func() {
			count := 3 + rand.Int()%7
			var newCharacter *character.Character
			for i := 0; i < count; i++ {
				newCharacter = new(character.Character)
				Expect(faker.FakeData(newCharacter)).To(Succeed())
				characters = append(characters, newCharacter)
			}
			pbs = characters.ToPb().Characters
		})
		AfterEach(func() {
			Expect(pbs).To(HaveLen(len(characters)))
			for idx, pb := range pbs {
				if pb == nil {
					continue
				}
				Expect(pb.CharacterId).To(Equal(characters[idx].Id.String()))
				Expect(pb.OwnerId).To(Equal(characters[idx].OwnerId))
				Expect(pb.Name).To(Equal(characters[idx].Name))
				Expect(pb.Gender).To(Equal(string(characters[idx].Gender)))
				Expect(pb.Realm).To(Equal(string(characters[idx].Realm)))
				Expect(pb.PlayTime).To(Equal(characters[idx].PlayTime))
				Expect(pb.Location).To(Equal(characters[idx].Location.ToPb()))
				Expect(pb.DimensionId).To(Equal(characters[idx].DimensionId))
			}
		})
	})

	Describe("character validation", func() {
		It("should return nil if the name is valid", func() {
			Expect(c.ValidateName()).To(Succeed())
		})
		Context("given invalid data", func() {
			When("given an invalid name", func() {
				It("should error if it is to short", func() {
					c.Name = string(faker.Username()[0])
					Expect(c.ValidateName()).NotTo(Succeed())
					c.Name = faker.Username()[:1]
					Expect(c.ValidateName()).NotTo(Succeed())
				})
				It("should error if it is to long", func() {
					c.Name = faker.Username() + faker.Username() + faker.Username() + faker.Username() + faker.Username()
					Expect(c.ValidateName()).NotTo(Succeed())
				})
				It("should error if it contains invalid characters", func() {
					c.Name = faker.Username() + "!"
					Expect(c.ValidateName()).NotTo(Succeed())
				})
				It("should error if it contains profane word", func() {
					c.Name = faker.Username()[:2] + "shit" + faker.Username()[:2]
					Expect(c.ValidateName()).NotTo(Succeed())
					c.Name = faker.Username()[:2] + "sh1t" + faker.Username()[:2]
					Expect(c.ValidateName()).NotTo(Succeed())
				})
				It("should error if it contains a hidden profane word", func() {
					c.Name = faker.Username()[:2] + "sh1t" + faker.Username()[:2]
				})
				AfterEach(func() {
					Expect(c.ValidateName()).NotTo(Succeed())
				})
			})
			When("given an invalid realm", func() {
				It("should error if it is unknown", func() {
					c.Realm = game.Realm(faker.Username())
					Expect(c.ValidateRealm()).NotTo(Succeed())
				})
			})
			When("given an invalid gender", func() {
				It("should error if it is unknown", func() {
					c.Gender = game.Gender(faker.Username())
					Expect(c.ValidateGender()).NotTo(Succeed())
				})
			})
			AfterEach(func() {
				Expect(c.Validate()).NotTo(Succeed())
			})
		})
	})
})
