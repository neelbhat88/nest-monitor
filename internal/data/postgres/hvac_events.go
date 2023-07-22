package postgres

import (
	"context"
	"neelbhat88/nest-monitor/m/v2/internal/models"
	"strings"

	"github.com/rs/zerolog/log"
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
