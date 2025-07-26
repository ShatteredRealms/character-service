package repository_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	crepository "github.com/ShatteredRealms/go-common-service/pkg/repository"
	"github.com/ShatteredRealms/go-common-service/pkg/testsro"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
)

type initializeData struct {
	PostgresConfig config.DBConfig
	MdbConnStr     string
	RedisConfig    config.DBPoolConfig
	DimensionId    uuid.UUID
}

var (
	_    = log.Logger
	hook *test.Hook

	pgCloseFunc  func() error
	mdb          *mongo.Database
	mdbCloseFunc func() error

	data     initializeData
	migrater *crepository.PgxMigrater
)

func TestRepository(t *testing.T) {
	var err error

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
	})
	SynchronizedBeforeSuite(func(ctx SpecContext) []byte {
		// log.Logger, hook = test.NewNullLogger()

		var pgPort string
		pgCloseFunc, pgPort, err = testsro.SetupPostgresWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(pgCloseFunc).NotTo(BeNil())

		mdbCloseFunc, data.MdbConnStr, err = testsro.SetupMongoWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(mdbCloseFunc).NotTo(BeNil())

		data.PostgresConfig = config.DBConfig{
			ServerAddress: config.ServerAddress{
				Port: pgPort,
				Host: "localhost",
			},
			Name:     testsro.DbName,
			Username: testsro.Username,
			Password: testsro.Password,
		}
		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())

		// migrater, err = crepository.NewPgxMigrater(ctx, data.PostgresConfig.PostgresDSN(), "migrations")
		// Expect(err).NotTo(HaveOccurred())
		// Expect(migrater).NotTo(BeNil())
		// ct, err := migrater.Conn.Exec(ctx, "CREATE TABLE dimensions (id TEXT PRIMARY KEY, updated_at TIMESTAMP, created_at TIMESTAMP);")
		// Expect(err).NotTo(HaveOccurred())
		// data.DimensionId, err = uuid.NewV7()
		// Expect(err).NotTo(HaveOccurred())
		// ct, err = migrater.Conn.Exec(ctx, "INSERT INTO dimensions (id, updated_at, created_at) VALUES (?, current_timestamp, current_timestamp);", data.DimensionId.String())
		// Expect(err).NotTo(HaveOccurred())
		// Expect(ct.RowsAffected()).To(Equal(1))

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		return buf.Bytes()
	}, func(ctx SpecContext, inBytes []byte) {
		// log.Logger, hook = test.NewNullLogger()

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())

		migrater, err = crepository.NewPgxMigrater(ctx, data.PostgresConfig.PostgresDSN(), "../../migrations", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(migrater).NotTo(BeNil())
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		pgCloseFunc()
		mdbCloseFunc()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
