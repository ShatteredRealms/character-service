package repository_test

import (
	"errors"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/ShatteredRealms/character-service/pkg/common"
	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	cgame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	cpb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
)

var _ = Describe("Character Repository", func() {
	var repo repository.CharacterRepository
	var c *character.Character

	BeforeEach(func() {
		var err error
		repo = repository.NewPgxCharacterRepository(migrater)
		Expect(err).NotTo(HaveOccurred())
		Expect(repo).NotTo(BeNil())
		c = RandomCharacter()
	})

	Describe("NewPostgresCharacter", func() {
		It("should be able to be called multiple times", func() {
			for i := 0; i < 2; i++ {
				repo = repository.NewPgxCharacterRepository(migrater)
			}
		})
	})

	Describe("CreateCharacter", func() {
		It("should require a non-nil character", func() {
			outC, err := repo.CreateCharacter(nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, common.ErrRequestInvalid)).To(BeTrue())
			Expect(outC).To(BeNil())
		})
		It("should require an empty id", func(ctx SpecContext) {
			uuid, err := uuid.NewV7()
			Expect(err).NotTo(HaveOccurred())

			c.Id = uuid
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, common.ErrRequestInvalid)).To(BeTrue())
			Expect(outC).To(BeNil())
		})
		It("should return a character", func(ctx SpecContext) {
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			ExpectCharactersEquals(outC, c)
		})
		It("should not allow duplicate names", func(ctx SpecContext) {
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			ExpectCharactersEquals(outC, c)

			outC, err = repo.CreateCharacter(ctx, c)
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())
		})
	})

	Context("with a character", func() {
		BeforeEach(func(ctx SpecContext) {
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			ExpectCharactersEquals(outC, c)

			c = outC

			outChars, _, err := repo.GetCharacters(ctx, map[string]interface{}{"id": c.Id}, nil, false)
			Expect(err).To(BeNil())
			Expect(outChars).To(ContainElement(c))
		})

		Describe("GetCharacterById", func() {
			It("should return nothing if there are no matches", func(ctx SpecContext) {
				outC, err := repo.GetCharacter(ctx, map[string]interface{}{"id": "a"})
				Expect(err).NotTo(BeNil())
				Expect(outC).To(BeNil())

				id, err := uuid.NewV7()
				Expect(err).NotTo(HaveOccurred())

				outC, err = repo.GetCharacter(ctx, map[string]interface{}{"id": id})
				Expect(err).To(BeNil())
				Expect(outC).To(BeNil())
			})
			It("should return a character if there is a match", func(ctx SpecContext) {
				outC, err := repo.GetCharacter(ctx, map[string]interface{}{"id": c.Id})
				Expect(err).To(BeNil())
				Expect(outC).NotTo(BeNil())
				ExpectCharactersEquals(outC, c)
			})
		})

		Describe("GetCharacters", func() {
			It("should return characters", func(ctx SpecContext) {
				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{}, &cpb.QueryFilters{}, false)
				Expect(err).To(BeNil())
				Expect(len(outChars) > 0).To(BeTrue())
				Expect(outChars).To(HaveLen(total))
			})

			It("should work with proto filter limit", func(ctx SpecContext) {
				c.Id = uuid.Nil
				c.Name = faker.Username() + "a"
				_, err := repo.CreateCharacter(ctx, c)
				Expect(err).To(BeNil())

				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{}, &cpb.QueryFilters{Limit: 1}, false)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(1))
				Expect(total > 1).To(BeTrue())
			})
			It("should work with proto filter offset", func(ctx SpecContext) {
				c.Id = uuid.Nil
				c.Name = faker.Username() + "a"
				_, err := repo.CreateCharacter(ctx, c)
				Expect(err).To(BeNil())

				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{}, &cpb.QueryFilters{Limit: 1}, false)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(1))
				Expect(total > 1).To(BeTrue())

				outChars2, total2, err := repo.GetCharacters(ctx, map[string]interface{}{}, &cpb.QueryFilters{Limit: 1, Offset: 1}, false)
				Expect(err).To(BeNil())
				Expect(outChars2).To(HaveLen(1))
				Expect(total2 > 1).To(BeTrue())
				Expect(outChars[0].Id).NotTo(Equal(outChars2[0].Id))
			})

			It("should work with proto match filters", func(ctx SpecContext) {
				request := &pb.Character{Id: c.Id.String(), Name: "a"}
				pbmask := &fieldmaskpb.FieldMask{Paths: []string{"id"}}
				mask, err := fieldmask_utils.MaskFromProtoFieldMask(pbmask, util.PascalCase)
				Expect(err).To(BeNil())

				matchFilters := make(map[string]interface{})
				Expect(fieldmask_utils.StructToMap(mask, request, matchFilters)).To(Succeed())

				outChars, _, err := repo.GetCharacters(ctx, matchFilters, &cpb.QueryFilters{}, false)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(1))
				Expect(outChars[0].Id).To(Equal(c.Id))
			})
		})

		Describe("GetCharactersByOwner", func() {
			It("should return nothing if there are no matches", func(ctx SpecContext) {
				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{"owner_id": "a"}, &cpb.QueryFilters{}, false)
				Expect(err).NotTo(BeNil())
				Expect(outChars).To(HaveLen(0))
				// Expect(outChars).To(HaveLen(total))

				id := uuid.New()
				outChars, total, err = repo.GetCharacters(ctx, map[string]interface{}{"owner_id": id.String()}, &cpb.QueryFilters{}, false)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(0))
				Expect(outChars).To(HaveLen(total))
			})

			It("should return a character if there is a match", func(ctx SpecContext) {
				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{"owner_id": c.OwnerId}, &cpb.QueryFilters{}, false)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(1))
				Expect(outChars).To(HaveLen(total))

				outC := outChars[0]
				Expect(err).To(BeNil())
				Expect(outC).NotTo(BeNil())
				ExpectCharactersEquals(outC, c)
			})
		})

		Describe("UpdateCharacter", func() {
			It("should error if the character is nil", func(ctx SpecContext) {
				outC, err := repo.UpdateCharacter(ctx, nil)
				Expect(err).To(HaveOccurred())
				Expect(outC).To(BeNil())
			})
			It("should update and return the character", func(ctx SpecContext) {
				c.Name = faker.Username() + "b"
				outC, err := repo.UpdateCharacter(ctx, c)
				Expect(err).To(BeNil())
				Expect(outC).NotTo(BeNil())
				ExpectCharactersEquals(outC, c)
			})
		})

		Describe("DeleteCharacter", func() {
			It("should delete nothing if there are no matches", func(ctx SpecContext) {
				id, err := uuid.NewV7()
				Expect(err).NotTo(HaveOccurred())
				outC, err := repo.DeleteCharacter(ctx, &id)
				Expect(err).To(BeNil())
				Expect(outC).To(BeNil())
			})

			It("should delete a character if there is a match", func(ctx SpecContext) {
				outC, err := repo.DeleteCharacter(ctx, &c.Id)
				Expect(err).To(BeNil())
				Expect(outC).NotTo(BeNil())
				ExpectCharactersEquals(outC, c)

				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{}, &cpb.QueryFilters{}, false)
				Expect(err).To(BeNil())
				Expect(outChars).NotTo(ContainElement(c))
				Expect(len(outChars) > 0).To(BeTrue())
				Expect(outChars).To(HaveLen(total))
			})
		})

		Describe("DeleteCharactersByOwner", func() {
			It("should delete nothing if there are no matches", func(ctx SpecContext) {
				id := uuid.New()
				outChars, err := repo.DeleteCharactersByOwner(ctx, &id)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(0))
			})

			It("should delete a character if there is a match", func(ctx SpecContext) {
				outChars, err := repo.DeleteCharactersByOwner(ctx, &c.OwnerId)
				Expect(err).To(BeNil())
				Expect(outChars).To(HaveLen(1))

				outC := outChars[0]
				Expect(outC).NotTo(BeNil())
				ExpectCharactersEquals(outC, c)

				outChars, total, err := repo.GetCharacters(ctx, map[string]interface{}{}, &cpb.QueryFilters{}, false)
				Expect(err).To(BeNil())
				Expect(outChars).NotTo(ContainElement(c))
				Expect(len(outChars) > 0).To(BeTrue())
				Expect(outChars).To(HaveLen(total))
			})
		})
	})
})

func ExpectCharactersEquals(a, b *character.Character) {
	Expect(a).NotTo(BeNil())
	Expect(a.Id).NotTo(BeNil())
	Expect(a.OwnerId).To(Equal(b.OwnerId))
	Expect(a.Name).To(Equal(b.Name))
	Expect(a.Gender).To(Equal(b.Gender))
	Expect(a.Realm).To(Equal(b.Realm))
	Expect(a.DimensionId).To(Equal(b.DimensionId))
	Expect(a.PlayTime).To(Equal(b.PlayTime))
	Expect(a.Location).To(Equal(b.Location))
}

func RandomCharacter() *character.Character {
	return &character.Character{
		OwnerId:     uuid.New(),
		Name:        faker.Username(),
		Gender:      game.GenderMale,
		Realm:       game.RealmHuman,
		DimensionId: uuid.New(),
		Location: cgame.Location{
			WorldId: uuid.New(),
			X:       rand.Float32(),
			Y:       rand.Float32(),
			Z:       rand.Float32(),
			Roll:    rand.Float32(),
			Pitch:   rand.Float32(),
			Yaw:     rand.Float32(),
		},
	}
}
