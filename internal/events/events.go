package events

import (
	"context"
	"neelbhat88/nest-monitor/m/v2/internal/data"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
)

func WriteMessage(ctx context.Context, dataController data.Controller, msg *pubsub.Message) error {
	log.Info().RawJSON("msg", msg.Data).Msg("Got message")

	err := dataController.WriteEventMessage(ctx, msg.Data)
	if err != nil {
		return err
	}

	msg.Ack()

	return nil
}
