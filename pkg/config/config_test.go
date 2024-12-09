package config_test

import (
	"fmt"
	"io"
	"os"

	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

var _ = Describe("NewCharacterConfig", func() {
	envKey := "SRO_SERVER_HOST"
	envVal := faker.Username()
	Expect(os.Setenv(envKey, envVal)).To(Succeed())

	var err error
	cfg, err := config.NewCharacterConfig(nil)
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())
	log.Logger.Out = io.Discard

	It("should return a new CharacterConfig", func() {
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

	It("should bind to environment variables", func() {
		Expect(cfg.Server.Host).To(Equal(envVal))
	})

	It("should read from a config file", func() {
		randStr := faker.Username() + "a"
		cfgData := []byte(fmt.Sprintf("keycloak:\n  realm: %s\nserver:\n  host: %s", randStr, randStr))
		Expect(os.WriteFile("sro-character.yaml", cfgData, 0644)).To(Succeed())
		defer os.Remove("sro-character.yaml")
		cfg2, err := config.NewCharacterConfig(nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg2).NotTo(BeNil())
		Expect(cfg2.Keycloak.Realm).To(Equal(randStr))
		Expect(cfg2.Server.Host).NotTo(Equal(randStr))
	})
})
