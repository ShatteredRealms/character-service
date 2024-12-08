package srv_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
)

var _ = Describe("CharacterContext", func() {
	Describe("NewCharacterContext", func() {
		var cfg *config.CharacterConfig
		BeforeEach(func(ctx SpecContext) {
			var err error
			cfg, err = config.NewCharacterConfig(ctx)
			Expect(err).To(BeNil())
			Expect(cfg).NotTo(BeNil())
			cfg.Postgres.Master = data.GormConfig
			cfg.Redis = data.RedisConfig
			cfg.Kafka = cconfig.ServerAddresses{data.KafkaConfig}
			Expect(cfg.Kafka).To(HaveLen(1))
		})
		Context("invalid input", func() {
			It("should error given invalid postgres connection", func() {
				cfg.Postgres.Master.Port = "0"
				cfg.Postgres.Master.Host = faker.Username()
			})
			It("should error given invalid redis connection", func() {
				cfg.Redis.Master.Port = "0"
				cfg.Redis.Master.Host = faker.Username()
				for idx := range cfg.Redis.Slaves {
					cfg.Redis.Slaves[idx].Port = "0"
					cfg.Redis.Slaves[idx].Host = faker.Username()
				}
			})

			AfterEach(func(ctx SpecContext) {
				cCtx, err := srv.NewCharacterContext(ctx, cfg, faker.Username())
				Expect(err).NotTo(BeNil())
				Expect(cCtx).To(BeNil())
			})
		})
		It("should create a valid context", func(ctx SpecContext) {
			var cCtx *srv.CharacterContext
			Eventually(func() error {
				var err error
				cCtx, err = srv.NewCharacterContext(ctx, cfg, faker.Username())
				return err
			}).Should(Succeed())
			Expect(cCtx).NotTo(BeNil())
		})
	})
})
