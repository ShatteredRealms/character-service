package repository_test

import (
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/repository"
	cgame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
)

var _ = Describe("CharacterPostgres", Ordered, func() {
	var repo repository.CharacterRepository
	var c, c2 *character.Character
	var dimensionId uuid.UUID

	BeforeAll(func() {
		var err error
		dimensionId = uuid.New()
		Expect(gdb.Exec(
			"CREATE TABLE dimensions (id TEXT PRIMARY KEY, updated_at TIMESTAMP, created_at TIMESTAMP);",
		).Error).NotTo(HaveOccurred())
		Expect(gdb.Exec(
			"INSERT INTO dimensions (id, updated_at, created_at) VALUES (?, current_timestamp, current_timestamp);",
			dimensionId,
		).Error).NotTo(HaveOccurred())
		repo, err = repository.NewPostgresCharacter(gdb)
		Expect(err).NotTo(HaveOccurred())
		Expect(repo).NotTo(BeNil())
		c = &character.Character{
			OwnerId:     uuid.New(),
			Name:        faker.Username(),
			Gender:      game.GenderMale,
			Realm:       game.RealmHuman,
			DimensionId: uuid.New(),
			PlayTime:    0,
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
	})

	Describe("NewPostgresCharacter", func() {
		It("should be able to be called multiple times", func() {
			for i := 0; i < 2; i++ {
				repo, err := repository.NewPostgresCharacter(gdb)
				Expect(err).NotTo(HaveOccurred())
				Expect(repo).NotTo(BeNil())
			}
		})
	})

	Describe("CreateCharacter", func() {
		It("should require a non-nil character", func() {
			// Expect(func() { repo.CreateCharacter(nil, nil) }).To(Panic())
		})
		It("should require an empty id", func(ctx SpecContext) {
			uuid, err := uuid.NewV7()
			Expect(err).NotTo(HaveOccurred())

			c.Id = &uuid
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())
		})
		It("should require a valid dimension to exist", func(ctx SpecContext) {
			c.Id = nil
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())
		})
		It("should return a character", func(ctx SpecContext) {
			c.DimensionId = dimensionId
			outC, err := repo.CreateCharacter(ctx, c)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC.Id).NotTo(BeNil())
			Expect(outC.OwnerId).To(Equal(c.OwnerId))
			Expect(outC.Name).To(Equal(c.Name))
			Expect(outC.Gender).To(Equal(c.Gender))
			Expect(outC.Realm).To(Equal(c.Realm))
			Expect(outC.DimensionId).To(Equal(c.DimensionId))
			Expect(outC.PlayTime).To(Equal(c.PlayTime))
			Expect(outC.Location).To(Equal(c.Location))
			c = outC
		})
		It("should not allow duplicate names", func(ctx SpecContext) {
			c2 = &character.Character{
				OwnerId:     uuid.New(),
				Name:        c.Name,
				Gender:      game.GenderMale,
				Realm:       game.RealmHuman,
				DimensionId: dimensionId,
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
			outC, err := repo.CreateCharacter(ctx, c2)
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())

			c2.Name = faker.Username() + "a"
			outC, err = repo.CreateCharacter(ctx, c2)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC.Id).NotTo(BeNil())
			Expect(outC.OwnerId).To(Equal(c2.OwnerId))
			Expect(outC.Name).To(Equal(c2.Name))
			Expect(outC.Gender).To(Equal(c2.Gender))
			Expect(outC.Realm).To(Equal(c2.Realm))
			Expect(outC.DimensionId).To(Equal(c2.DimensionId))
			Expect(outC.PlayTime).To(Equal(c2.PlayTime))
			Expect(outC.Location).To(Equal(c2.Location))
		})
	})

	Describe("GetCharacterById", func() {
		It("should return nothing if there are no matches", func(ctx SpecContext) {
			outC, err := repo.GetCharacterById(ctx, nil)
			Expect(err).To(BeNil())
			Expect(outC).To(BeNil())

			id, err := uuid.NewV7()
			Expect(err).NotTo(HaveOccurred())

			outC, err = repo.GetCharacterById(ctx, &id)
			Expect(err).To(BeNil())
			Expect(outC).To(BeNil())
		})
		It("should return a character if there is a match", func(ctx SpecContext) {
			outC, err := repo.GetCharacterById(ctx, c.Id)
			Expect(err).To(BeNil())
			Expect(outC).NotTo(BeNil())
			Expect(outC).NotTo(BeNil())
			Expect(outC.Id).To(Equal(c.Id))
			Expect(outC.OwnerId).To(Equal(c.OwnerId))
			Expect(outC.Name).To(Equal(c.Name))
			Expect(outC.Gender).To(Equal(c.Gender))
			Expect(outC.Realm).To(Equal(c.Realm))
			Expect(outC.DimensionId).To(Equal(c.DimensionId))
			Expect(outC.PlayTime).To(Equal(c.PlayTime))
			Expect(outC.Location).To(Equal(c.Location))
		})
	})

	Describe("GetCharacters", func() {
		It("should return characters", func() {
			outChars, err := repo.GetCharacters(nil)
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(2))
		})
	})

	Describe("GetCharactersByOwner", func() {
		It("should return nothing if there are no matches", func(ctx SpecContext) {
			outChars, err := repo.GetCharactersByOwner(ctx, "")
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(0))

			outChars, err = repo.GetCharactersByOwner(ctx, faker.UUIDHyphenated())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(0))
		})

		It("should return a character if there is a match", func(ctx SpecContext) {
			outChars, err := repo.GetCharactersByOwner(ctx, c.OwnerId.String())
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(1))

			outC := (*outChars)[0]
			Expect(err).To(BeNil())
			Expect(outC).NotTo(BeNil())
			Expect(outC).NotTo(BeNil())
			Expect(outC.Id).To(Equal(c.Id))
			Expect(outC.OwnerId).To(Equal(c.OwnerId))
			Expect(outC.Name).To(Equal(c.Name))
			Expect(outC.Gender).To(Equal(c.Gender))
			Expect(outC.Realm).To(Equal(c.Realm))
			Expect(outC.DimensionId).To(Equal(c.DimensionId))
			Expect(outC.PlayTime).To(Equal(c.PlayTime))
			Expect(outC.Location).To(Equal(c.Location))
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
			Expect(outC.Id).To(Equal(c.Id))
			Expect(outC.OwnerId).To(Equal(c.OwnerId))
			Expect(outC.Name).To(Equal(c.Name))
			Expect(outC.Gender).To(Equal(c.Gender))
			Expect(outC.Realm).To(Equal(c.Realm))
			Expect(outC.DimensionId).To(Equal(c.DimensionId))
			Expect(outC.PlayTime).To(Equal(c.PlayTime))
			Expect(outC.Location).To(Equal(c.Location))
		})
	})

	Describe("DeleteCharacter", func() {
		It("should delete nothing if there are no matches", func(ctx SpecContext) {
			id, err := uuid.NewV7()
			Expect(err).NotTo(HaveOccurred())
			outC, err := repo.DeleteCharacter(ctx, &id)
			Expect(err).To(BeNil())
			Expect(outC).To(BeNil())

			outChars, err := repo.GetCharacters(ctx)
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(2))
		})

		It("should delete a character if there is a match", func(ctx SpecContext) {
			outC, err := repo.DeleteCharacter(ctx, c.Id)
			Expect(err).To(BeNil())
			Expect(outC).NotTo(BeNil())
			Expect(outC.Id).To(Equal(c.Id))
			Expect(outC.OwnerId).To(Equal(c.OwnerId))
			Expect(outC.Name).To(Equal(c.Name))
			Expect(outC.Gender).To(Equal(c.Gender))
			Expect(outC.Realm).To(Equal(c.Realm))
			Expect(outC.DimensionId).To(Equal(c.DimensionId))
			Expect(outC.PlayTime).To(Equal(c.PlayTime))
			Expect(outC.Location).To(Equal(c.Location))

			outChars, err := repo.GetCharacters(ctx)
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(1))
		})
	})

	Describe("DeleteCharactersByOwner", func() {
		It("should delete nothing if there are no matches", func(ctx SpecContext) {
			outChars, err := repo.DeleteCharactersByOwner(ctx, faker.UUIDHyphenated())
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(0))

			outChars, err = repo.GetCharacters(ctx)
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(1))
		})

		It("should delete a character if there is a match", func(ctx SpecContext) {
			outChars, err := repo.DeleteCharactersByOwner(ctx, c2.OwnerId.String())
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(1))

			outC := (*outChars)[0]
			Expect(outC).NotTo(BeNil())
			Expect(outC.Id).To(Equal(c2.Id))
			Expect(outC.OwnerId).To(Equal(c2.OwnerId))
			Expect(outC.Name).To(Equal(c2.Name))
			Expect(outC.Gender).To(Equal(c2.Gender))
			Expect(outC.Realm).To(Equal(c2.Realm))
			Expect(outC.DimensionId).To(Equal(c2.DimensionId))
			Expect(outC.PlayTime).To(Equal(c2.PlayTime))
			Expect(outC.Location).To(Equal(c2.Location))

			outChars, err = repo.GetCharacters(ctx)
			Expect(err).To(BeNil())
			Expect(outChars).NotTo(BeNil())
			Expect(*outChars).To(HaveLen(0))
		})
	})
})
