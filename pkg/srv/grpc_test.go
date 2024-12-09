package srv_test

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/srv"
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
		BeforeEach(func(ctx SpecContext) {
			server, err := srv.NewCharacterServiceServer(ctx, cCtx)
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
	})
})
