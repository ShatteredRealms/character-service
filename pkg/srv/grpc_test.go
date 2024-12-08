package srv_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
)

var _ = Describe("Grpc", func() {
	var cfg *config.CharacterConfig
	var cCtx *srv.CharacterContext
	BeforeEach(func(ctx SpecContext) {
		var err error
		cfg, err = config.NewCharacterConfig(ctx)
		Expect(err).To(BeNil())
		Expect(cfg).NotTo(BeNil())
		cfg.Postgres.Master = data.GormConfig
		cfg.Redis = data.RedisConfig
		cfg.Kafka = cconfig.ServerAddresses{data.KafkaConfig}
		Expect(cfg.Kafka).To(HaveLen(1))
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
	Context("valid config", func() {
		BeforeEach(func(ctx SpecContext) {
			server, err := srv.NewCharacterServiceServer(ctx, cCtx)
			Expect(err).To(BeNil())
			Expect(server).NotTo(BeNil())
		})
		It("should allow creation of a new server", func() {
		})
	})
})
