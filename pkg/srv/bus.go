package srv

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

func (s *characterServiceServer) ListenMessageBus() {
	go s.listenDimension()
}

func (s *characterServiceServer) listenDimension() {
	for {
		func() {
			ctx := context.Background()
			msg, err := s.Context.DimensionBusReader.FetchMessage(ctx)

			ctx, span := s.Context.Tracer.Start(ctx, "sro.character.bus.dimension")
			defer span.End()

			if err != nil {
				log.Logger.WithContext(ctx).Errorf("dimension bus: %v", err)
				return
			}

			log.Logger.WithContext(ctx).Infof("parsing dimension bus id: %v, deleted: %v", msg.Id, msg.Deleted)

			if msg.Deleted {
				dimension, err := s.Context.DimensionService.DeleteDimension(ctx, msg.Id)
				if err != nil {
					log.Logger.WithContext(ctx).Errorf("parse dimension '%s' failed: %v", msg.Id, err)
					s.Context.DimensionBusReader.ProcessFailed()
					return
				}
				if dimension == nil {
					log.Logger.WithContext(ctx).Warnf("dimension '%s' did not exist, parsing complete", msg.Id)
				} else {
					log.Logger.WithContext(ctx).Infof("parse dimension '%s' succeeded", msg.Id)
				}
				s.Context.DimensionBusReader.ProcessSucceeded(ctx)
				return
			}

			_, err = s.Context.DimensionService.CreateDimension(ctx, msg.Id)
			if err != nil {
				log.Logger.WithContext(ctx).Errorf("parse dimension '%s' failed: %v", msg.Id, err)
				s.Context.DimensionBusReader.ProcessFailed()
				return
			}

			log.Logger.WithContext(ctx).Infof("created dimension: %s", msg.Id)
			s.Context.DimensionBusReader.ProcessSucceeded(ctx)
		}()
	}
}
