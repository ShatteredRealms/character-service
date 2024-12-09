package srv_test

import (
	"context"
	"testing"
	"time"

	"bytes"
	"encoding/gob"

	"github.com/WilSimpson/gocloak/v13"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"

	"github.com/ShatteredRealms/character-service/pkg/config"
	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/testsro"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type initializeData struct {
	GormConfig     cconfig.DBConfig
	MdbConnStr     string
	RedisConfig    cconfig.DBPoolConfig
	KafkaConfig    cconfig.ServerAddress
	KeycloakConfig cconfig.KeycloakConfig
}

var (
	hook *test.Hook

	gdb               *gorm.DB
	gdbCloseFunc      func() error
	mdb               *mongo.Database
	mdbCloseFunc      func() error
	redisCloseFunc    func() error
	kafkaCloseFunc    func() error
	keycloakCloseFunc func() error

	data initializeData
	cfg  *config.CharacterConfig

	admin = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testadmin"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(false),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("adminfirstname"),
		LastName:      gocloak.StringP("adminlastname"),
		Email:         gocloak.StringP("admin@example.com"),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}
	user = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testplayer"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(false),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("userfirstname"),
		LastName:      gocloak.StringP("userlastname"),
		Email:         gocloak.StringP("user@example.com"),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}

	incCtxAdmin context.Context
	incCtxUser  context.Context
)

func TestSrv(t *testing.T) {
	var err error
	createCfg := func() {
		var err error
		cfg, err = config.NewCharacterConfig(nil)
		Expect(err).To(BeNil())
		Expect(cfg).NotTo(BeNil())
	}
	setupCfgFunc := func() {
		cfg.Postgres.Master = data.GormConfig
		cfg.Redis = data.RedisConfig
		cfg.Kafka = cconfig.ServerAddresses{data.KafkaConfig}
		cfg.Keycloak.BaseURL = data.KeycloakConfig.BaseURL
		Expect(cfg.Kafka).To(HaveLen(1))
	}

	SynchronizedBeforeSuite(func() []byte {
		log.Logger, hook = test.NewNullLogger()
		createCfg()

		data.KeycloakConfig = cfg.Keycloak
		keycloakCloseFunc, data.KeycloakConfig.BaseURL, err = testsro.SetupKeycloakWithDocker()
		Expect(err).To(BeNil())

		var gormPort string
		gdbCloseFunc, gormPort, err = testsro.SetupGormWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(gdbCloseFunc).NotTo(BeNil())

		mdbCloseFunc, data.MdbConnStr, err = testsro.SetupMongoWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(mdbCloseFunc).NotTo(BeNil())

		redisCloseFunc, data.RedisConfig, err = testsro.SetupRedisWithDocker()
		Expect(err).To(BeNil())

		data.KafkaConfig = cfg.Kafka[0]
		kafkaCloseFunc, data.KafkaConfig.Port, err = testsro.SetupKafkaWithDocker()
		Expect(err).To(BeNil())

		data.GormConfig = cconfig.DBConfig{
			ServerAddress: cconfig.ServerAddress{
				Port: gormPort,
				Host: "localhost",
			},
			Name:     testsro.DbName,
			Username: testsro.Username,
			Password: testsro.Password,
		}
		gdb, err = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(err).NotTo(HaveOccurred())
		Expect(gdb).NotTo(BeNil())
		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		setupCfgFunc()

		keycloak := gocloak.NewClient(data.KeycloakConfig.BaseURL)
		Expect(keycloak).NotTo(BeNil())

		clientToken, err := keycloak.LoginClient(
			context.Background(),
			cfg.Keycloak.ClientId,
			cfg.Keycloak.ClientSecret,
			cfg.Keycloak.Realm,
		)
		Expect(err).NotTo(HaveOccurred())

		*admin.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, cfg.Keycloak.Realm, admin)
		Expect(err).NotTo(HaveOccurred())
		*user.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, cfg.Keycloak.Realm, user)
		Expect(err).NotTo(HaveOccurred())

		saRole, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, cfg.Keycloak.Realm, "super admin")
		Expect(err).NotTo(HaveOccurred())
		userRole, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, cfg.Keycloak.Realm, "user")
		Expect(err).NotTo(HaveOccurred())

		err = keycloak.AddRealmRoleToUser(
			context.Background(),
			clientToken.AccessToken,
			cfg.Keycloak.Realm,
			*admin.ID,
			[]gocloak.Role{*saRole},
		)
		Expect(err).NotTo(HaveOccurred())
		err = keycloak.AddRealmRoleToUser(
			context.Background(),
			clientToken.AccessToken,
			cfg.Keycloak.Realm,
			*user.ID,
			[]gocloak.Role{*userRole},
		)
		Expect(err).NotTo(HaveOccurred())

		return buf.Bytes()
	}, func(inBytes []byte) {
		log.Logger, hook = test.NewNullLogger()

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())
		createCfg()
		setupCfgFunc()

		gdb, err = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(err).NotTo(HaveOccurred())
		Expect(gdb).NotTo(BeNil())
		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())

		keycloak := gocloak.NewClient(cfg.Keycloak.BaseURL)
		Expect(keycloak).NotTo(BeNil())

		clientToken, err := keycloak.LoginClient(
			context.Background(),
			cfg.Keycloak.ClientId,
			cfg.Keycloak.ClientSecret,
			cfg.Keycloak.Realm,
		)
		Expect(err).NotTo(HaveOccurred())

		setupUser := func(user *gocloak.User, ctx *context.Context) {
			var token *gocloak.JWT
			Eventually(func() error {
				token, err = keycloak.Login(
					context.Background(),
					cfg.Keycloak.ClientId,
					cfg.Keycloak.ClientSecret,
					cfg.Keycloak.Realm,
					*user.Username,
					*(*user.Credentials)[0].Value,
				)
				return err
			}).Within(time.Minute).Should(Succeed())
			Expect(err).NotTo(HaveOccurred())

			users, err := keycloak.GetUsers(
				context.Background(),
				clientToken.AccessToken,
				cfg.Keycloak.Realm,
				gocloak.GetUsersParams{Username: user.Username},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(users).To(HaveLen(1))
			user = users[0]

			md := metadata.New(
				map[string]string{
					"authorization": "Bearer " + token.AccessToken,
				},
			)
			(*ctx) = metadata.NewIncomingContext(context.Background(), md)
		}

		// setupUser(&admin, &incCtxAdmin)
		setupUser(&user, &incCtxUser)

	})

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		createCfg()
		setupCfgFunc()
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		if gdbCloseFunc != nil {
			gdbCloseFunc()
		}
		if mdbCloseFunc != nil {
			mdbCloseFunc()
		}
		if redisCloseFunc != nil {
			redisCloseFunc()
		}
		if kafkaCloseFunc != nil {
			kafkaCloseFunc()
		}
		if keycloakCloseFunc != nil {
			keycloakCloseFunc()
		}
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")

	It("should work", func() {
	})
}
