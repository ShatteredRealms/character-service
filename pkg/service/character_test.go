package service_test

import (
	"errors"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	mock_repository "github.com/ShatteredRealms/character-service/pkg/repository/mocks"
	"github.com/ShatteredRealms/character-service/pkg/service"
	cgame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
)

var _ = Describe("CharacterS", func() {
	var repo *mock_repository.MockCharacterRepository
	var svc service.CharacterService
	var c *character.Character
	BeforeEach(func() {
		controller := gomock.NewController(GinkgoT())
		Expect(controller).NotTo(BeNil())
		repo = mock_repository.NewMockCharacterRepository(controller)
		Expect(repo).NotTo(BeNil())
		svc = service.NewCharacterService(repo)
		Expect(svc).NotTo(BeNil())
		dimensionId := uuid.New()
		c = &character.Character{
			OwnerId:     uuid.New(),
			Name:        faker.Username(),
			Gender:      game.GenderMale,
			Realm:       game.RealmHuman,
			DimensionId: dimensionId,
			PlayTime:    rand.Uint64(),
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

	Describe("AddCharacterPlaytime", func() {
		It("should fail if updating the repo fails", func(ctx SpecContext) {
			repo.EXPECT().UpdateCharacter(gomock.Any(), gomock.Any()).Return(nil, errors.New("repo"))
			outC, err := svc.AddCharacterPlaytime(ctx, c, rand.Uint64())
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())
		})
		It("should update the playtime if updating the player succeeds", func(ctx SpecContext) {
			timeAdded := 300 + rand.Uint64()%700
			originalTime := c.PlayTime
			repo.EXPECT().UpdateCharacter(gomock.Any(), gomock.Any()).Return(c, nil)
			outC, err := svc.AddCharacterPlaytime(ctx, c, timeAdded)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC).To(Equal(c))
			Expect(outC.PlayTime).To(Equal(originalTime + timeAdded))
		})
	})

	Describe("CreateCharacter", func() {
		Context("failure", func() {
			It("should fail if the character name is to short", func() {
				c.Name = faker.Username()[:0]
			})
			It("should fail if the character name is to short", func() {
				c.Name = faker.Username()[:1]
			})
			It("should fail if the character name is to long", func() {
				c.Name = faker.Username() + faker.Username() + faker.Username() + faker.Username()
			})
			It("should fail if the character name is profane", func() {
				c.Name = "sh1t"
			})
			It("should fail if the repo creation fails", func() {
				repo.EXPECT().CreateCharacter(gomock.Any(), gomock.Any()).Return(nil, errors.New("repo"))
			})
			AfterEach(func(ctx SpecContext) {
				outC, err := svc.CreateCharacter(ctx, c.OwnerId.String(), c.Name, string(c.Gender), string(c.Realm), &c.DimensionId)
				Expect(err).To(HaveOccurred())
				Expect(outC).To(BeNil())
			})
		})
		It("should create a character if it is valid", func(ctx SpecContext) {
			repo.EXPECT().CreateCharacter(gomock.Eq(ctx), gomock.Any()).Return(c, nil)
			outC, err := svc.CreateCharacter(ctx, c.OwnerId.String(), c.Name, string(c.Gender), string(c.Realm), &c.DimensionId)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC).To(Equal(c))
		})
	})
	Describe("EditCharacter", func() {
		Context("failure", func() {
			It("should fail if the character name is to short", func() {
				c.Name = faker.Username()[:0]
			})
			It("should fail if the character name is to short", func() {
				c.Name = faker.Username()[:1]
			})
			It("should fail if the character name is to long", func() {
				c.Name = faker.Username() + faker.Username() + faker.Username() + faker.Username()
			})
			It("should fail if the character name is profane", func() {
				c.Name = "sh1t"
			})
			It("should fail if the repo creation fails", func() {
				repo.EXPECT().UpdateCharacter(gomock.Any(), gomock.Eq(c)).Return(nil, errors.New("repo"))
			})
			AfterEach(func(ctx SpecContext) {
				outC, err := svc.EditCharacter(ctx, c)
				Expect(err).To(HaveOccurred())
				Expect(outC).To(BeNil())
			})
		})
		It("should create a character if it is valid", func(ctx SpecContext) {
			repo.EXPECT().CreateCharacter(gomock.Eq(ctx), gomock.Any()).Return(c, nil)
			outC, err := svc.CreateCharacter(ctx, c.OwnerId.String(), c.Name, string(c.Gender), string(c.Realm), &c.DimensionId)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC).To(Equal(c))
		})
	})
	Describe("DeleteCharacter", func() {
		It("should fail if the repo fails", func(ctx SpecContext) {
			repo.EXPECT().DeleteCharacter(gomock.Eq(ctx), gomock.Eq(&c.Id)).Return(nil, errors.New("repo"))
			outC, err := svc.DeleteCharacter(ctx, &c.Id)
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())
		})
		It("should delete the player if the repo succeeds", func(ctx SpecContext) {
			repo.EXPECT().DeleteCharacter(gomock.Eq(ctx), gomock.Eq(&c.Id)).Return(c, nil)
			outC, err := svc.DeleteCharacter(ctx, &c.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC).To(Equal(c))
		})
	})
	Describe("GetCharacterById", func() {
		It("should fail if the repo fails", func(ctx SpecContext) {
			repo.EXPECT().GetCharacterById(gomock.Eq(ctx), gomock.Eq(&c.Id)).Return(nil, errors.New("repo"))
			outC, err := svc.GetCharacterById(ctx, &c.Id)
			Expect(err).To(HaveOccurred())
			Expect(outC).To(BeNil())
		})
		It("should return the results of the repo if it succeeds", func(ctx SpecContext) {
			repo.EXPECT().GetCharacterById(gomock.Eq(ctx), gomock.Eq(&c.Id)).Return(c, nil)
			outC, err := svc.GetCharacterById(ctx, &c.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(outC).NotTo(BeNil())
			Expect(outC).To(Equal(c))
		})
	})
	Describe("GetCharacters", func() {
		It("should fail if the character does not exist", func(ctx SpecContext) {
			repo.EXPECT().GetCharacters(gomock.Eq(ctx)).Return(nil, errors.New("repo"))
			outChars, err := svc.GetCharacters(ctx)
			Expect(err).To(HaveOccurred())
			Expect(outChars).To(BeNil())
		})
		It("should return the results of the repo if it succeeds", func(ctx SpecContext) {
			repo.EXPECT().GetCharacters(gomock.Eq(ctx)).Return(character.Characters{c}, nil)
			outChars, err := svc.GetCharacters(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(outChars).NotTo(BeNil())
			Expect(outChars).To(HaveLen(1))
			Expect(outChars[0]).To(Equal(c))
		})
	})
	Describe("GetCharactersByOwner", func() {
		It("should fail if the character does not exist", func(ctx SpecContext) {
			repo.EXPECT().GetCharactersByOwner(gomock.Eq(ctx), gomock.Eq(&c.OwnerId)).Return(nil, errors.New("repo"))
			outChars, err := svc.GetCharactersByOwner(ctx, c.OwnerId.String())
			Expect(err).To(HaveOccurred())
			Expect(outChars).To(BeNil())
		})
		It("should return the results of the repo if it succeeds", func(ctx SpecContext) {
			repo.EXPECT().GetCharactersByOwner(gomock.Eq(ctx), gomock.Eq(&c.OwnerId)).Return(character.Characters{c}, nil)
			outChars, err := svc.GetCharactersByOwner(ctx, c.OwnerId.String())
			Expect(err).NotTo(HaveOccurred())
			Expect(outChars).NotTo(BeNil())
			Expect(outChars).To(HaveLen(1))
			Expect(outChars[0]).To(Equal(c))
		})
	})
})
