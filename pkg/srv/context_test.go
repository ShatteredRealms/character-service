package srv_test

import (
	"errors"
	"math/rand/v2"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CharacterContext", func() {
	Describe("Close", func() {
		It("should close the dimension service", func() {
			dimensionServiceMock.EXPECT().StopProcessing()
			srvCtx.Close()
			Expect(srvCtx.DimensionService).NotTo(BeNil())
		})
		It("should not close the dimension service if it is nil", func() {
			srvCtx.DimensionService = nil
			srvCtx.Close()
			Expect(srvCtx.DimensionService).To(BeNil())
		})
	})

	Describe("ResetCharacterBus", func() {
		var fn commonsrv.WriterResetCallback

		BeforeEach(func() {
			fn = srvCtx.ResetCharacterBus()
		})

		Context("generating errors", func() {
			var outErr error
			expectedErr := errors.New(faker.Username())

			It("should return an error if getting characters fails", func(ctx SpecContext) {
				mockCharService.EXPECT().GetCharacters(gomock.Any()).Return(nil, -1, expectedErr)
				outErr = fn(ctx)
			})

			When("getting characters succeeded", func() {
				BeforeEach(func() {
					mockCharService.EXPECT().GetCharacters(gomock.Any()).Return([]*character.Character{NewCharacter()}, 1, nil)
				})
				It("should error if getting deleted characters fails", func(ctx SpecContext) {
					mockCharService.EXPECT().GetDeletedCharacters(gomock.Any()).Return(nil, -1, expectedErr)
					outErr = fn(ctx)
				})
				When("getting deleted characters succeeds", func() {
					BeforeEach(func() {
						mockCharService.EXPECT().GetDeletedCharacters(gomock.Any()).Return([]*character.Character{NewCharacter()}, 1, nil)
					})
					It("should return error when writing to the bus fails", func(ctx SpecContext) {
						characterBusWriterMock.EXPECT().PublishMany(gomock.Any(), gomock.Any()).Return(expectedErr)
						outErr = fn(ctx)
					})
				})
			})

			AfterEach(func() {
				Expect(errors.Is(outErr, expectedErr)).To(BeTrue())
			})
		})
		It("should return nil if all operations succeed", func(ctx SpecContext) {
			existingChars := make([]*character.Character, rand.IntN(5)+2)
			for idx := 0; idx < len(existingChars); idx++ {
				existingChars[idx] = NewCharacter()
				GinkgoWriter.Printf("Existing character: %v\n", existingChars[idx].Id)
			}
			deletedChars := make([]*character.Character, rand.IntN(5)+2)
			for idx := 0; idx < len(deletedChars); idx++ {
				deletedChars[idx] = NewCharacter()
				GinkgoWriter.Printf("Deleted character: %v\n", deletedChars[idx].Id)
			}

			mockCharService.EXPECT().GetCharacters(gomock.Any()).Return(existingChars, len(existingChars), nil)
			mockCharService.EXPECT().GetDeletedCharacters(gomock.Any()).Return(deletedChars, len(existingChars), nil)
			characterBusWriterMock.EXPECT().
				PublishMany(gomock.Any(), gomock.Any()).
				Return(nil).
				Do(func(ctx any, msgs []characterbus.Message) {
					Expect(len(msgs)).To(Equal(len(existingChars) + len(deletedChars)))
				})
			Expect(fn(ctx)).To(BeNil())
		})
	})
})
