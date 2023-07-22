package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"neelbhat88/nest-monitor/m/v2/internal/models"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetWeather(ctx context.Context, apiKey string) (models.Weather, []byte, error) {
	url := fmt.Sprintf("https://api.tomorrow.io/v4/weather/realtime?location=60067&units=imperial&apikey=%v", apiKey)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error calling Tomorrow API")
		return models.Weather{}, []byte{}, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var weather models.Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshaling Weather data")
		return models.Weather{}, []byte{}, err
	}

	return weather, body, nil
}
