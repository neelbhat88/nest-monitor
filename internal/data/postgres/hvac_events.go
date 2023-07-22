package postgres

import (
	"context"
	"neelbhat88/nest-monitor/m/v2/internal/models"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	STATUS_NOISE = "noise"
)

func (n NestMonitorDB) WriteHvacEvent(ctx context.Context, event models.HVACEvent) error {
	_, err := n.Exec(`
		INSERT INTO hvac_events(event_id, event_timestamp, hvac_status)
		VALUES($1, $2, $3)
	`, event.EventID, event.Timestamp, strings.ToLower(event.ResourceUpdate.Traits.ThermostatHVAC.Status))
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", event.EventID).
			Str("timestamp", event.Timestamp).
			Str("hvac_status", event.ResourceUpdate.Traits.ThermostatHVAC.Status).
			Msg("Insert into hvac_events failed")

		return err
	}

	return nil
}

func (n NestMonitorDB) WriteNoiseEvent(ctx context.Context) (string, error) {
	var id string
	err := n.QueryRowx(`
		INSERT INTO hvac_events(event_timestamp, hvac_status)
		VALUES($1, $2)
		RETURNING event_id
	`, time.Now().UTC(), STATUS_NOISE).Scan(&id)
	if err != nil {
		log.Error().
			Err(err).
			Str("hvac_status", STATUS_NOISE).
			Msg("Insert noise event into hvac_events failed")

		return "", err
	}

	return id, nil
}
