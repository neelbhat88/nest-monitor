package postgres

import (
	"context"
	"neelbhat88/nest-monitor/m/v2/internal/models"

	"github.com/rs/zerolog/log"
)

func (n NestMonitorDB) WriteWeather(ctx context.Context, eventID string, weather models.Weather, rawData []byte) error {
	_, err := n.Exec(`
		INSERT INTO weather(event_id, temperature_apparent_f, temperature_f, humidity, pressure_surface_level, wind_speed, raw_data)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`, eventID, weather.Data.Values.TemperatureApparent, weather.Data.Values.Temperature, weather.Data.Values.Humidity, weather.Data.Values.PressureSurfaceLevel, weather.Data.Values.WindSpeed, rawData)
	if err != nil {
		log.Error().Err(err).Any("weather", weather).Msg("Insert into weather failed")
		return err
	}

	return nil
}
