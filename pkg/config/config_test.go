package config_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

var _ = Describe("NewCharacterConfig", func() {
	It("should return a new CharacterConfig", func() {
		log.Logger.Out = io.Discard
		cfg, err := config.NewCharacterConfig(nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())
		Expect(cfg.Server.Host).NotTo(BeEmpty())
		Expect(cfg.Server.Port).NotTo(BeEmpty())
		Expect(cfg.Keycloak.BaseURL).To(ContainSubstring("http://"))
		Expect(cfg.Keycloak.BaseURL).To(ContainSubstring("http://"))
		Expect(cfg.Keycloak.Realm).NotTo(BeEmpty())
		Expect(cfg.Keycloak.Id).NotTo(BeEmpty())
		Expect(cfg.Keycloak.ClientId).NotTo(BeEmpty())
		Expect(cfg.Keycloak.ClientSecret).NotTo(BeEmpty())
		Expect(cfg.Mode).NotTo(BeEmpty())
		Expect(cfg.OpenTelemtryAddress).NotTo(BeEmpty())
		Expect(cfg.Kafka).NotTo(BeEmpty())
		Expect(cfg.Postgres.Master.Name).NotTo(ContainSubstring("-"))
		Expect(cfg.Postgres.Master.Host).NotTo(BeEmpty())
		Expect(cfg.Postgres.Master.Port).NotTo(BeEmpty())
		Expect(cfg.Postgres.Master.Name).NotTo(BeEmpty())
		Expect(cfg.Postgres.Master.Username).NotTo(BeEmpty())
		Expect(cfg.Postgres.Master.Password).NotTo(BeEmpty())
		Expect(cfg.Redis.Master.Host).NotTo(BeEmpty())
		Expect(cfg.Redis.Master.Port).NotTo(BeEmpty())
		Expect(cfg.Redis.Master.Port).NotTo(BeEmpty())
	})
})
