package srv_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/WilSimpson/gocloak/v13"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	mock_service "github.com/ShatteredRealms/character-service/pkg/service/mocks"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	mock_dimensionbus "github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus/mocks"
	"github.com/ShatteredRealms/go-common-service/pkg/model"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	commonpb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

var _ = Describe("Grpc Server", func() {
	var cCtx *srv.CharacterContext
	Describe("Setting up the server", func() {
		BeforeEach(func(ctx SpecContext) {
			var err error
			Eventually(func() error {
				cCtx, err = srv.NewCharacterContext(ctx, cfg, faker.Username())
				return err
			}).Should(Succeed())
			Expect(cCtx).NotTo(BeNil())
			Expect(cCtx.DimensionService.IsProcessing()).To(BeTrue())
		})

		It("should fail to create given an invalid keycloak connection", func(ctx SpecContext) {
			cCtx.Config.Keycloak.Realm = faker.Username()
			server, err := srv.NewCharacterServiceServer(ctx, cCtx)
			Expect(err).NotTo(BeNil())
			Expect(server).To(BeNil())
		})

		It("should have roles", func() {
			Expect(srv.CharacterRoles).NotTo(BeEmpty())
		})
	})

	Context("using the character service", func() {
		var server pb.CharacterServiceServer
		var dim *dimensionbus.Dimension
		var userChar *character.Character
		var adminChar *character.Character

		var controller *gomock.Controller
		var mockCharService *mock_service.MockCharacterService
		var mockDimService *mock_dimensionbus.MockService
		var busWriter MockCharacterBusWriter

		BeforeEach(func(ctx SpecContext) {
			controller = gomock.NewController(GinkgoT())
			mockCharService = mock_service.NewMockCharacterService(controller)
			mockDimService = mock_dimensionbus.NewMockService(controller)
			cCtx = &srv.CharacterContext{
				Context: &commonsrv.Context{
					Config:         &cfg.BaseConfig,
					KeycloakClient: gocloak.NewClient(cfg.Keycloak.BaseURL),
					Tracer:         otel.Tracer(fmt.Sprintf("test-%d", GinkgoParallelProcess())),
				},
				CharacterBusWriter: busWriter,
				CharacterService:   mockCharService,
				DimensionService:   mockDimService,
			}

			var err error
			server, err = srv.NewCharacterServiceServer(ctx, cCtx)
			Expect(err).To(BeNil())
			Expect(server).NotTo(BeNil())

			id, err := uuid.NewV7()
			Expect(err).To(BeNil())
			dim = &dimensionbus.Dimension{
				Id:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			id, err = uuid.NewV7()
			Expect(err).To(BeNil())

			userChar = &character.Character{
				Model: model.Model{
					Id:        &id,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				OwnerId:     uuid.MustParse(*user.ID),
				Name:        faker.Username(),
				Gender:      game.GenderMale,
				Realm:       game.RealmHuman,
				DimensionId: dim.Id,
				Dimension:   dim,
				PlayTime:    rand.Uint64N(1000),
				Location: commongame.Location{
					WorldId: uuid.New(),
				},
			}

			id, err = uuid.NewV7()
			Expect(err).To(BeNil())
			adminChar = &character.Character{
				Model: model.Model{
					Id:        &id,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				OwnerId:     uuid.MustParse(*admin.ID),
				Name:        faker.Username(),
				Gender:      game.GenderMale,
				Realm:       game.RealmHuman,
				DimensionId: dim.Id,
				Dimension:   dim,
				PlayTime:    rand.Uint64N(1000),
				Location: commongame.Location{
					WorldId: uuid.New(),
				},
			}
		})
		Describe("NewCharacterServiceServer", func() {
			It("should allow creation of a new server", func() {
			})

			for _, role := range srv.CharacterRoles {
				It(fmt.Sprintf("should have role %s created", *role.Name), func(ctx SpecContext) {
					srvJwt, err := cCtx.GetJWT(ctx)
					Expect(err).To(BeNil())
					outRole, err := cCtx.KeycloakClient.GetClientRole(
						ctx,
						srvJwt.AccessToken,
						cfg.Keycloak.Realm,
						cfg.Keycloak.Id,
						*role.Name,
					)
					Expect(err).To(BeNil())
					Expect(outRole).NotTo(BeNil())
					Expect(*outRole.Name).To(Equal(*role.Name))
				})
			}
		})
		Context("valid permissions", func() {
			Describe("AddCharacterPlayTime", func() {
				It("should add play time to a character", func(ctx SpecContext) {
					amount := rand.Uint64N(1000)
					mockCharService.EXPECT().GetCharacterById(inCtxAdmin, userChar.Id).Return(userChar, nil)
					mockCharService.EXPECT().AddCharacterPlaytime(inCtxAdmin, userChar, amount).Return(userChar, nil)
					out, err := server.AddCharacterPlayTime(inCtxAdmin, &pb.AddPlayTimeRequest{
						CharacterId: userChar.Id.String(),
						Time:        amount,
					})
					Expect(err).To(BeNil())
					Expect(out).NotTo(BeNil())
				})
				It("should fail if adding to playtime fails", func(ctx SpecContext) {
					amount := rand.Uint64N(1000)
					retErr := errors.New(faker.Username())
					mockCharService.EXPECT().GetCharacterById(inCtxAdmin, userChar.Id).Return(userChar, nil)
					mockCharService.EXPECT().AddCharacterPlaytime(inCtxAdmin, userChar, amount).Return(nil, retErr)
					out, err := server.AddCharacterPlayTime(inCtxAdmin, &pb.AddPlayTimeRequest{
						CharacterId: userChar.Id.String(),
						Time:        amount,
					})
					Expect(err).NotTo(BeNil())

					sts, ok := status.FromError(err)
					Expect(ok).To(BeTrue(), "expected error to be a status error")
					Expect(sts.Code()).To(Equal(codes.Internal), "expected internal but got %s", sts.Message())
					Expect(sts.Message()).To(ContainSubstring(srv.ErrCharacterPlaytime.Error()))
					Expect(out).To(BeNil())
				})
				It("should fail if the given character id is invalid", func(ctx SpecContext) {
					out, err := server.AddCharacterPlayTime(inCtxAdmin, &pb.AddPlayTimeRequest{})
					Expect(err).NotTo(BeNil())

					sts, ok := status.FromError(err)
					Expect(ok).To(BeTrue(), "expected error to be a status error")
					Expect(sts.Code()).To(Equal(codes.InvalidArgument), "expected invalid argument but got %s", sts.Message())
					Expect(sts.Message()).To(ContainSubstring(srv.ErrCharacterIdInvalid.Error()))
					Expect(out).To(BeNil())
				})
				It("should fail if the given character does not exist", func(ctx SpecContext) {
					mockCharService.EXPECT().GetCharacterById(inCtxAdmin, userChar.Id).Return(nil, nil)
					out, err := server.AddCharacterPlayTime(inCtxAdmin, &pb.AddPlayTimeRequest{
						CharacterId: userChar.Id.String(),
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
		Context("invalid permissions", func() {
			for _, call := range []ServerCall{
				{
					function: "AddCharacterPlayTime",
					fn: func(ctx context.Context) (any, error) {
						return server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
							CharacterId: userChar.Id.String(),
							Time:        100 + rand.Uint64N(900),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user"},
					},
				},

				{
					function: "CreateCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.CreateCharacter(ctx, &pb.CreateCharacterRequest{})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
					},
				},
				{
					function: "CreateCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.CreateCharacter(ctx, &pb.CreateCharacterRequest{
							OwnerId: *admin.ID,
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user"},
					},
				},

				{
					function: "DeleteCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.DeleteCharacter(ctx, &commonpb.TargetId{
							Id: adminChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user",
							extraFn: func() { mockCharService.EXPECT().GetCharacterById(inCtxUser, adminChar.Id).Return(adminChar, nil) }},
					},
				},
				{
					function: "DeleteCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.DeleteCharacter(ctx, &commonpb.TargetId{
							Id: userChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
					},
				},

				{
					function: "EditCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.EditCharacter(ctx, &pb.EditCharacterRequest{
							CharacterId: adminChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user"},
					},
				},
				{
					function: "EditCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.EditCharacter(ctx, &pb.EditCharacterRequest{
							CharacterId: userChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
					},
				},

				{
					function: "GetCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharacter(ctx, &commonpb.TargetId{
							Id: userChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
					},
				},
				{
					function: "GetCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharacter(ctx, &commonpb.TargetId{
							Id: adminChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user",
							extraFn: func() { mockCharService.EXPECT().GetCharacterById(inCtxUser, adminChar.Id).Return(adminChar, nil) }},
					},
				},

				{
					function: "GetCharacters",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharacters(ctx, &emptypb.Empty{})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user"},
					},
				},

				{
					function: "GetCharactersForUser (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharactersForUser(ctx, &commonpb.TargetId{
							Id: *admin.ID,
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
					},
				},
				{
					function: "GetCharactersForUser (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharactersForUser(ctx, &commonpb.TargetId{
							Id: adminChar.Id.String(),
						})
					},
					calls: []ServerCallDetail{
						{ctx: nil, name: "guest"},
						{ctx: &inCtxUser, name: "user"},
					},
				},
			} {
				Describe(call.function, func() {
					for _, details := range call.calls {
						It(fmt.Sprintf("should fail for role %s", details.name), func(ctx context.Context) {
							if details.ctx == nil {
								details.ctx = &ctx
							}
							if details.extraFn != nil {
								details.extraFn()
							}
							out, err := call.fn(*details.ctx)
							Expect(out).To(BeNil())
							Expect(err).To(Equal(commonsrv.ErrPermissionDenied))
						})
					}
				})
			}
		})
	})
})

type ServerFunc func(ctx context.Context) (any, error)
type ServerCallDetail struct {
	ctx     *context.Context
	name    string
	extraFn func()
}
type ServerCall struct {
	function string
	fn       ServerFunc
	calls    []ServerCallDetail
}

type MockCharacterBusWriter struct {
	RetErr error
}

func (m MockCharacterBusWriter) Publish(ctx context.Context, message characterbus.Message) error {
	return m.RetErr
}

func (m MockCharacterBusWriter) Close() error {
	return m.RetErr
}

func (m MockCharacterBusWriter) GetMessageType() bus.BusMessageType {
	return "test"
}

func (m MockCharacterBusWriter) PublishMany(ctx context.Context, messages []characterbus.Message) error {
	return m.RetErr
}
