package srv_test

import (
	"io"
	"testing"

	mock_service "github.com/ShatteredRealms/character-service/pkg/service/mocks"
	"github.com/ShatteredRealms/character-service/pkg/srv"
	"github.com/ShatteredRealms/go-common-service/pkg/auth"
	"github.com/ShatteredRealms/go-common-service/pkg/bus/character/characterbus"
	mock_dimensionbus "github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/dimensionbus/mocks"
	mock_bus "github.com/ShatteredRealms/go-common-service/pkg/bus/mocks"
	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	commonsrv "github.com/ShatteredRealms/go-common-service/pkg/srv"
	mock_gocloak "github.com/WilSimpson/gocloak/v13/mocks"
	iauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
)

var (
	srvCtx                 *srv.CharacterContext
	mockConrol             *gomock.Controller
	dimensionServiceMock   *mock_dimensionbus.MockService
	mockCharService        *mock_service.MockCharacterService
	characterBusWriterMock *mock_bus.MockMessageBusWriter[characterbus.Message]
	keycloakClientMock     *mock_gocloak.MockGoCloakIface
	authFunc               iauth.AuthFunc
)

const (
	kcServiceName = "character"
)

func TestSrv(t *testing.T) {
	BeforeEach(func() {
		log.Logger.Out = io.Discard
		mockConrol = gomock.NewController(GinkgoT())
		keycloakClientMock = mock_gocloak.NewMockGoCloakIface(mockConrol)
		dimensionServiceMock = mock_dimensionbus.NewMockService(mockConrol)
		mockCharService = mock_service.NewMockCharacterService(mockConrol)
		characterBusWriterMock = mock_bus.NewMockMessageBusWriter[characterbus.Message](mockConrol)
		srvCtx = &srv.CharacterContext{
			Context: &commonsrv.Context{
				Config: &config.BaseConfig{
					Keycloak: config.KeycloakConfig{
						ClientId: kcServiceName,
					},
				},
				KeycloakClient: keycloakClientMock,
				Tracer:         otel.Tracer("TestCharacterContext"),
			},
			CharacterBusWriter: characterBusWriterMock,
			CharacterService:   mockCharService,
			DimensionService:   dimensionServiceMock,
		}
		authFunc = auth.AuthFunc(keycloakClientMock, "realm")
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")
}
