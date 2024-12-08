package srv_test

import (
	"testing"

	"bytes"
	"encoding/gob"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/testsro"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type initializeData struct {
	GormConfig  config.DBConfig
	MdbConnStr  string
	RedisConfig config.DBPoolConfig
	KafkaConfig config.ServerAddress
}

var (
	hook *test.Hook

	gdb            *gorm.DB
	gdbCloseFunc   func() error
	mdb            *mongo.Database
	mdbCloseFunc   func() error
	redisCloseFunc func() error
	kafkaCloseFunc func() error

	data initializeData
)

func TestSrv(t *testing.T) {
	var err error
	SynchronizedBeforeSuite(func() []byte {
		log.Logger, hook = test.NewNullLogger()

		var gormPort string
		gdbCloseFunc, gormPort, err = testsro.SetupGormWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(gdbCloseFunc).NotTo(BeNil())

		mdbCloseFunc, data.MdbConnStr, err = testsro.SetupMongoWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(mdbCloseFunc).NotTo(BeNil())

		redisCloseFunc, data.RedisConfig, err = testsro.SetupRedisWithDocker()
		Expect(err).To(BeNil())

		data.KafkaConfig.Host = "localhost"
		kafkaCloseFunc, data.KafkaConfig.Port, err = testsro.SetupKafkaWithDocker()

		data.GormConfig = config.DBConfig{
			ServerAddress: config.ServerAddress{
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

		return buf.Bytes()
	}, func(inBytes []byte) {
		log.Logger, hook = test.NewNullLogger()

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		gdb, err = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(err).NotTo(HaveOccurred())
		Expect(gdb).NotTo(BeNil())
		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())
	})

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
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
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")
}
