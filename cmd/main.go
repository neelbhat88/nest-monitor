package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"neelbhat88/nest-monitor/m/v2/internal/data/postgres"
	"neelbhat88/nest-monitor/m/v2/internal/events"
	"neelbhat88/nest-monitor/m/v2/internal/weather"

	"cloud.google.com/go/pubsub"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type AppConfig struct {
	GCloudProjectID      string `env:"GCLOUD_PROJECTID"`
	GCloudSubscriptionID string `env:"GCLOUD_SUBSCRIPTIONID"`
	ADCFile              string `env:"ADC_FILE"`
	TomorrowAPIKey       string `env:"TOMORROW_API_KEY"`
}

type DatabaseMigrationSource struct {
}

func (DatabaseMigrationSource) GetMigrations() postgres.PostgresMigrations {
	return postgres.GetMigrations()
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("No .env file found")
	}

	var dbConfig postgres.DatabaseConfig
	err = cleanenv.ReadEnv(&dbConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read DatabaseConfig from config.env")
	}

	var appConfig AppConfig
	err = cleanenv.ReadEnv(&appConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read AppConfig from config.env")
	}

	ms := DatabaseMigrationSource{}
	db, err := postgres.InitializeDB(dbConfig, ms)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize DB")
	}
	defer db.Close()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	nestRepo := postgres.NewPostgresRepository(db)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var id string
		err := nestRepo.DB.Get(&id, `
			select id
			from event_stream
			order by created desc
			limit 1
		`)
		if err != nil {
			log.Error().Err(err).Msg("Error reading from DB")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}

		w.Write([]byte(fmt.Sprintf("hi %v", id)))
	})

	r.Get("/noise", func(w http.ResponseWriter, r *http.Request) {
		id, err := nestRepo.WriteNoiseEvent(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Error writing Noise event")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}

		// Get current weather data and log it
		weather, rawBody, err := weather.GetWeather(r.Context(), appConfig.TomorrowAPIKey)
		if err != nil {
			log.Error().Err(err).Msg("Error getting weather data")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting weather data"))
			return
		}

		err = nestRepo.WriteWeather(r.Context(), id, weather, rawBody)
		if err != nil {
			log.Error().Err(err).Msg("Error writing weather data")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error writing weather data"))
			return
		}

		w.Write([]byte(fmt.Sprintf("Noise recorded! ID: %v", id)))
	})

	r.Get("/weather", func(w http.ResponseWriter, r *http.Request) {
		_, rawBody, err := weather.GetWeather(r.Context(), appConfig.TomorrowAPIKey)
		if err != nil {
			log.Error().Err(err).Msg("Error getting weather data")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting weather data"))
			return
		}

		w.Write(rawBody)
	})

	log.Info().Msg("Application started")
	log.Info().Str("path", appConfig.ADCFile).Msg("ADCFile")

	ctx := context.Background()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Stack().Msg("Panic Recovered")
			}
		}()

		client, err := pubsub.NewClient(ctx, appConfig.GCloudProjectID, option.WithCredentialsFile(appConfig.ADCFile))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create a new secretmanager client")
		}
		defer client.Close()

		sub := client.Subscription(appConfig.GCloudSubscriptionID)
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			err := events.WriteMessage(ctx, nestRepo, m)
			if err != nil {
				log.Error().Err(err).Msg("Error writing message to DB")
			}
		})
		if err != nil {
			log.Error().Err(err).Msg("Error on receiving message")
		}
	}()

	http.ListenAndServe(":8080", r)
}
