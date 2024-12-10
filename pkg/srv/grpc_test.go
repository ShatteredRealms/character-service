package srv_test

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
)

var _ = Describe("Grpc", func() {
	var cCtx *srv.CharacterContext
	BeforeEach(func(ctx SpecContext) {
		var err error
		Eventually(func() error {
			cCtx, err = srv.NewCharacterContext(ctx, cfg, faker.Username())
			return err
		}).Should(Succeed())
		Expect(cCtx).NotTo(BeNil())
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
	Context("valid config", func() {
		var server pb.CharacterServiceServer
		BeforeEach(func(ctx SpecContext) {
			var err error
			server, err = srv.NewCharacterServiceServer(ctx, cCtx)
			Expect(err).To(BeNil())
			Expect(server).NotTo(BeNil())
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
		Context("invalid permissions", func() {
			for _, call := range []ServerCall{
				{
					fn: func(ctx context.Context) (any, error) {
						return server.AddCharacterPlayTime(ctx, &pb.AddPlayTimeRequest{
							CharacterId: faker.UUIDHyphenated(),
							Time:        uint64(100 + rand.UintN(900)),
						})
					},
					function: "AddCharacterPlayTime",
					calls: []ServerCallDetail{
						{
							ctx:  &inCtxUser,
							name: "user",
						},
						{
							ctx:  nil,
							name: "guest",
						},
					},
				},
				{
					fn: func(ctx context.Context) (any, error) {
						return server.CreateCharacter(ctx, &pb.CreateCharacterRequest{})
					},
					function: "CreateCharacter (self)",
					calls: []ServerCallDetail{
						{
							ctx:  nil,
							name: "guest",
						},
					},
				},
				{
					fn: func(ctx context.Context) (any, error) {
						return server.CreateCharacter(ctx, &pb.CreateCharacterRequest{
							OwnerId: *admin.ID,
						})
					},
					function: "CreateCharacter (other)",
					calls: []ServerCallDetail{
						{
							ctx:  nil,
							name: "guest",
						},
						{
							ctx:  &inCtxUser,
							name: "user",
						},
					},
				},
			} {
				Describe(call.function, func() {
					for _, details := range call.calls {
						It(fmt.Sprintf("should fail for role %s", details.name), func(ctx context.Context) {
							GinkgoWriter.Printf("in inCtxAdmin: %v\n", inCtxAdmin)
							GinkgoWriter.Printf("in inCtxUser: %v\n", inCtxUser)
							if details.ctx == nil {
								details.ctx = &ctx
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
	ctx  *context.Context
	name string
}
type ServerCall struct {
	function string
	fn       ServerFunc
	calls    []ServerCallDetail
}
