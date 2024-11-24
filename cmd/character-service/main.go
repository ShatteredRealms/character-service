package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ShatteredRealms/character-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/srv"
	"github.com/ShatteredRealms/go-common-service/pkg/telemetry"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/WilSimpson/gocloak/v13"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srvErr := make(chan error, 1)
	go func() {
		defer func() {
			srvErr <- nil
		}()

		cfg, err := config.NewCharacterConfig(ctx)
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("loading config: %v", err)
			return
		}

		srvCtx := srv.NewContext(&cfg.BaseConfig, "CharacterService")
		ctx, span := srvCtx.Tracer.Start(ctx, "main")
		defer span.End()

		otelShutdown, err := telemetry.SetupOTelSDK(ctx, "character", config.Version, cfg.OpenTelemtryAddress)
		defer func() {
			log.Logger.Infof("Shutting down")
			err = otelShutdown(context.Background())
			if err != nil {
				log.Logger.Warnf("Error shutting down: %v", err)
			}
		}()

		if err != nil {
			log.Logger.WithContext(ctx).Errorf("connecting to otel: %v", err)
			return
		}

		span.End()
		srvErr <- util.StartServer(ctx, grpcServer, gwmux, cfg.Server.Address())
	}()

	select {
	case err := <-srvErr:
		if err != nil {
			log.Logger.Error(err)
		}

	case <-ctx.Done():
		log.Logger.Info("Server canceled by user input.")
		stop()
	}

	log.Logger.Info("Server stopped.")
}

func SetupHealthServer(ctx context.Context, cfg *config.BaseConfig) error {
	keycloakClient := gocloak.NewClient(cfg.Keycloak.BaseURL)
	grpcServer, gwmux := util.InitServerDefaults(keycloakClient, cfg.Keycloak.Realm)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, cfg.Server.Address(), opts)
	if err != nil {
		return fmt.Errorf("register health service handler endpoint: %w", err)
	}
	return nil
}
