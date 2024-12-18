package srv_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/WilSimpson/gocloak/v13"
	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ShatteredRealms/character-service/pkg/model/character"
	"github.com/ShatteredRealms/character-service/pkg/model/game"
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	"github.com/ShatteredRealms/go-common-service/pkg/auth"
	commongame "github.com/ShatteredRealms/go-common-service/pkg/model/game"
	commonpb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

var _ = Describe("Grpc Server", func() {
	Describe("Setting up the server", func() {
		It("should fail if creating roles fails", func(ctx SpecContext) {
			err := errors.New(faker.Username())
			keycloakClientMock.EXPECT().
				LoginClient(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&gocloak.JWT{}, nil)
			keycloakClientMock.EXPECT().
				CreateClientRole(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).
				Times(len(srv.CharacterRoles)).
				Return("", &gocloak.APIError{})
			server, err := srv.NewCharacterServiceServer(ctx, srvCtx)
			Expect(err).To(Equal(err))
			Expect(server).To(BeNil())
		})

		It("should setup roles for creation", func() {
			Expect(srv.CharacterRoles).NotTo(BeEmpty())
		})

		It("should create roles", func(ctx SpecContext) {
			keycloakClientMock.EXPECT().
				LoginClient(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&gocloak.JWT{}, nil).Times(2)
			keycloakClientMock.EXPECT().
				CreateClientRole(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).
				Times(len(srv.CharacterRoles)+len(srv.CompositeCharacterRoles)).
				Return("", nil)
			server, err := srv.NewCharacterServiceServer(ctx, srvCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(server).NotTo(BeNil())
		})
	})

	Context("using the character service", func() {
		var server pb.CharacterServiceServer
		var userChar *character.Character
		var adminChar *character.Character

		BeforeEach(func(ctx SpecContext) {
			var err error
			keycloakClientMock.EXPECT().
				LoginClient(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&gocloak.JWT{}, nil).AnyTimes()
			keycloakClientMock.EXPECT().
				CreateClientRole(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).
				Return("", nil).
				AnyTimes()
			server, err = srv.NewCharacterServiceServer(ctx, srvCtx)
			Expect(err).To(BeNil())
			Expect(server).NotTo(BeNil())

			userChar = NewCharacter()
			adminChar = NewCharacter()
		})
		Context("valid permissions", func() {
			Describe("AddCharacterPlayTime", func() {
				It("should add play time to a character", func() {
					amount := rand.Uint64N(1000)
					ctx := ReturnClaimsWithRole(srv.RolePlaytime, userChar.OwnerId.String())
					mockCharService.EXPECT().GetCharacterById(ctx, &userChar.Id).Return(userChar, nil)
					mockCharService.EXPECT().AddCharacterPlaytime(ctx, userChar, amount).Return(userChar, nil)
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
						CharacterId: userChar.Id.String(),
						Time:        amount,
					})
					Expect(err).To(BeNil())
					Expect(out).NotTo(BeNil())
				})
				It("should fail if adding to playtime fails", func() {
					amount := rand.Uint64N(1000)
					retErr := errors.New(faker.Username())
					ctx := ReturnClaimsWithRole(srv.RolePlaytime, userChar.OwnerId.String())
					mockCharService.EXPECT().GetCharacterById(ctx, &userChar.Id).Return(userChar, nil)
					mockCharService.EXPECT().AddCharacterPlaytime(ctx, userChar, amount).Return(nil, retErr)
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
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
				It("should fail if the given character id is invalid", func() {
					ctx := ReturnClaimsWithRole(srv.RolePlaytime, userChar.OwnerId.String())
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{})
					Expect(err).NotTo(BeNil())

					sts, ok := status.FromError(err)
					Expect(ok).To(BeTrue(), "expected error to be a status error")
					Expect(sts.Code()).To(Equal(codes.InvalidArgument), "expected invalid argument but got %s", sts.Message())
					Expect(sts.Message()).To(ContainSubstring(srv.ErrCharacterIdInvalid.Error()))
					Expect(out).To(BeNil())
				})
				It("should fail if the given character does not exist", func() {
					ctx := ReturnClaimsWithRole(srv.RolePlaytime, userChar.OwnerId.String())
					mockCharService.EXPECT().GetCharacterById(ctx, &userChar.Id).Return(nil, nil)
					out, err := server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
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

			Describe("CreateCharacter", func() {
				It("should create a character for self", func() {
					// ctx := ReturnClaimsWithRole(srv.RoleCharacterManagement, userChar.OwnerId.String())
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
					roles: []*gocloak.Role{
						srv.RolePlaytime,
					},
				},

				{
					function: "CreateCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.CreateCharacter(ctx, &pb.CreateCharacterRequest{})
					},
					roles: []*gocloak.Role{
						srv.RoleCreateCharactersSelf,
						srv.RoleCreateCharactersAll,
					},
				},
				{
					function: "CreateCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.CreateCharacter(ctx, &pb.CreateCharacterRequest{
							OwnerId: adminChar.OwnerId.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleCreateCharactersAll,
					},
				},

				{
					function: "DeleteCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						mockCharService.EXPECT().GetCharacterById(ctx, &adminChar.Id).Return(adminChar, nil).AnyTimes()
						return server.DeleteCharacter(ctx, &commonpb.TargetId{
							Id: adminChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleDeleteCharactersAll,
					},
				},
				{
					function: "DeleteCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.DeleteCharacter(ctx, &commonpb.TargetId{
							Id: userChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleDeleteCharactersSelf,
						srv.RoleDeleteCharactersAll,
					},
				},

				{
					function: "EditCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.EditCharacter(ctx, &pb.EditCharacterRequest{
							CharacterId: adminChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleEditCharacter,
					},
				},
				{
					function: "EditCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.EditCharacter(ctx, &pb.EditCharacterRequest{
							CharacterId: userChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleEditCharacter,
					},
				},

				{
					function: "GetCharacter (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharacter(ctx, &commonpb.TargetId{
							Id: userChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleGetCharactersSelf,
						srv.RoleGetCharactersAll,
					},
				},
				{
					function: "GetCharacter (other)",
					fn: func(ctx context.Context) (any, error) {
						mockCharService.EXPECT().GetCharacterById(ctx, &adminChar.Id).Return(adminChar, nil).AnyTimes()
						return server.GetCharacter(ctx, &commonpb.TargetId{
							Id: adminChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleGetCharactersAll,
					},
				},

				{
					function: "GetCharacters",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharacters(ctx, &emptypb.Empty{})
					},
					roles: []*gocloak.Role{
						srv.RoleGetCharactersAll,
					},
				},

				{
					function: "GetCharactersForUser (owner)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharactersForUser(ctx, &commonpb.TargetId{
							Id: userChar.OwnerId.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleGetCharactersSelf,
						srv.RoleGetCharactersAll,
					},
				},
				{
					function: "GetCharactersForUser (other)",
					fn: func(ctx context.Context) (any, error) {
						return server.GetCharactersForUser(ctx, &commonpb.TargetId{
							Id: adminChar.Id.String(),
						})
					},
					roles: []*gocloak.Role{
						srv.RoleGetCharactersAll,
					},
				},
			} {
				Describe(call.function, func() {
					for _, role := range srv.CharacterRoles {
						match := false
						for _, ignoreRole := range call.roles {
							if role.Name == ignoreRole.Name {
								match = true
								break
							}
						}
						if !match {
							It(fmt.Sprintf("should fail for role %s", *role.Name), func() {
								ctx := ReturnClaimsWithRole(role, userChar.OwnerId.String())
								out, err := call.fn(ctx)
								Expect(out).To(BeNil())
								Expect(err).To(Equal(commonsrv.ErrPermissionDenied))
							})
						}
					}
				})
			}
		})
	})
})

func ReturnClaimsWithRole(role *gocloak.Role, subject string) context.Context {
	claims := auth.SROClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour),
			},
			NotBefore: &jwt.NumericDate{
				Time: time.Now().Add(-time.Minute),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now().Add(-time.Hour),
			},
		},
		ResourceAccess: auth.ClaimResourceAccess{
			kcServiceName: auth.ClaimRoles{
				Roles: []string{*role.Name},
			},
		},
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"authorization": []string{fmt.Sprintf("Bearer %s", faker.UUIDHyphenated())},
	})

	keycloakClientMock.EXPECT().DecodeAccessTokenCustomClaims(
		gomock.Eq(ctx),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).SetArg(3, claims).Return(&jwt.Token{Valid: true}, nil)

	ctx, err := authFunc(ctx)
	Expect(err).To(BeNil())

	return ctx
}

func NewCharacter() *character.Character {
	return &character.Character{
		Id:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerId:     uuid.New(),
		Name:        faker.Username(),
		Gender:      game.GenderMale,
		Realm:       game.RealmHuman,
		DimensionId: uuid.New(),
		PlayTime:    rand.Uint64N(1000),
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
