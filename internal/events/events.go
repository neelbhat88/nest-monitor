package events

import (
	"context"
	"encoding/json"
	"neelbhat88/nest-monitor/m/v2/internal/data"
	"neelbhat88/nest-monitor/m/v2/internal/models"

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

	hvacEvent, err := UnmarshalMessage(ctx, msg.Data)
	if err != nil {
		return err
	}

	// Ignore any events without an HVAC status since those are not hvac events we care about
	if hvacEvent.ResourceUpdate.Traits.ThermostatHVAC.Status == "" {
		return nil
	}

	err = dataController.WriteHvacEvent(ctx, hvacEvent)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalMessage(ctx context.Context, msg []byte) (models.HVACEvent, error) {
	var event models.HVACEvent

	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Unmarshal message")
		return event, err
	}

	return event, nil
}
