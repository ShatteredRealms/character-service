package srv_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bytes"
	"encoding/gob"

	"github.com/WilSimpson/gocloak/v13"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/auth"
	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus"
	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
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
	ServerToken    string
	AdminToken     string
	UserToken      string
}

var (
	hook *test.Hook

	gdb *gorm.DB
	mdb *mongo.Database

	cfg *config.CharacterConfig

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

	inCtxServer context.Context
	inCtxAdmin  context.Context
	inCtxUser   context.Context
)

func TestSrv(t *testing.T) {
	var (
		err               error
		gdbCloseFunc      func() error
		mdbCloseFunc      func() error
		redisCloseFunc    func() error
		kafkaCloseFunc    func() error
		keycloakCloseFunc func() error
	)

	BeforeEach(func() {
		// log.Logger, hook = test.NewNullLogger()
		GinkgoWriter.Printf("Process Parallel Proccess: %d\n", GinkgoParallelProcess())
		GinkgoWriter.Printf("Postgres Config: %+v\n", cfg.Postgres)
	})

	SynchronizedBeforeSuite(func() []byte {
		var err error

		cfg, err = config.NewCharacterConfig(nil)
		Expect(err).To(BeNil())
		Expect(cfg).NotTo(BeNil())

		var data initializeData

		data.KeycloakConfig = cfg.Keycloak
		keycloakCloseFunc, data.KeycloakConfig.BaseURL, err = testsro.SetupKeycloakWithDocker()
		Expect(err).To(BeNil())

		data.GormConfig = cconfig.DBConfig{
			ServerAddress: cconfig.ServerAddress{
				Host: "localhost",
			},
			Name:     testsro.DbName,
			Username: testsro.Username,
			Password: testsro.Password,
		}
		gdbCloseFunc, data.GormConfig.Port, err = testsro.SetupPostgresWithDocker()
		fmt.Printf("Gorm Config: %+v\n", data.GormConfig)
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

		writer := bus.NewKafkaMessageBusWriter(cconfig.ServerAddresses{data.KafkaConfig}, dimensionbus.Message{})
		Eventually(func() error {
			return writer.Publish(context.Background(), dimensionbus.Message{
				Id:      uuid.New(),
				Deleted: false,
			})
		}).Within(time.Minute).Should(Succeed())
		Expect(writer.Close()).To(Succeed())

		keycloak := gocloak.NewClient(data.KeycloakConfig.BaseURL)
		Expect(keycloak).NotTo(BeNil())

		clientToken, err := keycloak.LoginClient(
			context.Background(),
			cfg.Keycloak.ClientId,
			cfg.Keycloak.ClientSecret,
			cfg.Keycloak.Realm,
		)
		Expect(err).NotTo(HaveOccurred())
		data.ServerToken = clientToken.AccessToken

		setupUser := func(user *gocloak.User, roleName string, tokenStr *string) {
			*user.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, cfg.Keycloak.Realm, *user)
			Expect(err).NotTo(HaveOccurred())
			role, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, cfg.Keycloak.Realm, roleName)
			Expect(err).NotTo(HaveOccurred())
			err = keycloak.AddRealmRoleToUser(
				context.Background(),
				clientToken.AccessToken,
				cfg.Keycloak.Realm,
				*user.ID,
				[]gocloak.Role{*role},
			)
			Expect(err).NotTo(HaveOccurred())
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
			(*tokenStr) = token.AccessToken
		}

		setupUser(&admin, "super admin", &data.AdminToken)
		setupUser(&user, "user", &data.UserToken)

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		return buf.Bytes()
	}, func(inBytes []byte) {
		var data initializeData

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		fmt.Printf("Gorm Config: %+v\n", data.GormConfig)
		gdb, err = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(err).NotTo(HaveOccurred())
		Expect(gdb).NotTo(BeNil())

		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())

		fn := auth.AuthFunc(gocloak.NewClient(data.KeycloakConfig.BaseURL), data.KeycloakConfig.Realm)

		inCtxServer, err = fn(metadata.NewIncomingContext(context.Background(), mdFn(data.ServerToken)))
		Expect(err).To(BeNil())
		inCtxAdmin, err = fn(metadata.NewIncomingContext(context.Background(), mdFn(data.AdminToken)))
		Expect(err).To(BeNil())
		inCtxUser, err = fn(metadata.NewIncomingContext(context.Background(), mdFn(data.UserToken)))
		Expect(err).To(BeNil())

		cfg, err = config.NewCharacterConfig(nil)
		Expect(err).To(BeNil())
		Expect(cfg).NotTo(BeNil())
		cfg.Postgres.Master = data.GormConfig
		cfg.Redis = data.RedisConfig
		cfg.Kafka = cconfig.ServerAddresses{data.KafkaConfig}
		cfg.Keycloak.BaseURL = data.KeycloakConfig.BaseURL
		Expect(cfg.Kafka).To(HaveLen(1))
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

func mdFn(token string) metadata.MD {
	return metadata.MD{
		"authorization": []string{"Bearer " + token},
	}
}
