package srv_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/WilSimpson/gocloak/v13"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/gender"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/realm"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
)

var _ = Describe("Grpc Server", func() {
	Describe("Setting up the server", func() {
		It("should allow creating a new server", func(ctx SpecContext) {
			server := srv.NewCharacterServiceServer(ctx, srvCtx)
			Expect(server).NotTo(BeNil(), "expected server to be created")
		})
	})

	Context("using the character service", func() {
		var server pb.CharacterServiceServer
		var userChar *character.Character

		BeforeEach(func(ctx SpecContext) {
			server = srv.NewCharacterServiceServer(ctx, srvCtx)
			Expect(server).NotTo(BeNil(), "expected server to be created")

			userChar = NewCharacter()
		})
		Context("valid permissions", func() {
			Describe("AddCharacterPlayTime", func() {
				It("should add play time to a character", func() {
					amount := rand.Int32N(1000)
					ctx := context.Background()
					mockCharService.EXPECT().GetCharacterById(ctx, &userChar.Id).Return(userChar, nil)
					mockCharService.EXPECT().AddCharacterPlaytime(ctx, userChar, amount).Return(userChar, nil)
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
						Id:   userChar.Id.String(),
						Time: amount,
					})
					Expect(err).To(BeNil())
					Expect(out).NotTo(BeNil())
				})
				It("should fail if adding to playtime fails", func() {
					amount := rand.Int32N(1000)
					retErr := errors.New(faker.Username())
					ctx := context.Background()
					mockCharService.EXPECT().GetCharacterById(ctx, &userChar.Id).Return(userChar, nil)
					mockCharService.EXPECT().AddCharacterPlaytime(ctx, userChar, amount).Return(nil, retErr)
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
						Id:   userChar.Id.String(),
						Time: amount,
					})
					Expect(err).NotTo(BeNil())

					sts, ok := status.FromError(err)
					Expect(ok).To(BeTrue(), "expected error to be a status error")
					Expect(sts.Code()).To(Equal(codes.Internal), "expected internal but got %s", sts.Message())
					Expect(sts.Message()).To(ContainSubstring(srv.ErrCharacterPlaytime.Error()))
					Expect(out).To(BeNil())
				})
				It("should fail if the given character id is invalid", func() {
					ctx := context.Background()
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{})
					Expect(err).NotTo(BeNil())

					sts, ok := status.FromError(err)
					Expect(ok).To(BeTrue(), "expected error to be a status error")
					Expect(sts.Code()).To(Equal(codes.InvalidArgument), "expected invalid argument but got %s", sts.Message())
					Expect(sts.Message()).To(ContainSubstring(srv.ErrCharacterIdInvalid.Error()))
					Expect(out).To(BeNil())
				})
				It("should fail if the given character does not exist", func() {
					ctx := context.Background()
					mockCharService.EXPECT().GetCharacterById(ctx, &userChar.Id).Return(nil, nil)
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
						Id: userChar.Id.String(),
					})
					Expect(err).NotTo(BeNil())

					sts, ok := status.FromError(err)
					Expect(ok).To(BeTrue(), "expected error to be a status error")
					Expect(sts.Code()).To(Equal(codes.NotFound), "expected not found but got %s", sts.Message())
					Expect(sts.Message()).To(ContainSubstring(srv.ErrCharacterDoesNotExist.Error()))
					Expect(out).To(BeNil())
				})
			})
		})
	})
})

func NewCharacter() *character.Character {
	return &character.Character{
		Id:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerId:     uuid.New(),
		Name:        faker.Username(),
		Gender:      gender.Male,
		Realm:       realm.Human,
		DimensionId: uuid.New(),
		PlayTime:    rand.Int32N(1000),
		Location: commongame.Location{
			WorldId: uuid.New(),
		},
	}
}

type ServerFunc func(ctx context.Context) (any, error)
type ServerCall struct {
	function string
	fn       ServerFunc
	roles    []*gocloak.Role
}
