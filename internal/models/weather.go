package models

type Weather struct {
	Data WeatherData `json:"data"`
}

type WeatherData struct {
	Values WeatherValues `json:"values"`
}

type WeatherValues struct {
	Humidity             int64   `json:"humidity"`
	PressureSurfaceLevel float64 `json:"pressureSurfaceLevel"`
	Temperature          float64 `json:"temperature"`
	TemperatureApparent  float64 `json:"temperatureApparent"`
	WindSpeed            float64 `json:"windSpeed"`
}
