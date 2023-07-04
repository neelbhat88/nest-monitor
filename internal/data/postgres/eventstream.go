package postgres

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (n NestMonitorDB) WriteEventMessage(ctx context.Context, msg []byte) error {
	_, err := n.Exec(`
		INSERT INTO event_stream(event_message)
		VALUES($1)
	`, msg)
	if err != nil {
		log.Error().Err(err).Msg("Insert into event_stream failed")
	}

	return nil
}
