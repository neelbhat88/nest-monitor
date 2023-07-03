package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Port     string `env:"PGPORT" env-default:"5432"`
	Host     string `env:"PGHOST" env-default:"localhost"`
	Name     string `env:"PGNAME" env-default:"postgres"`
	User     string `env:"PGUSER" env-default:"user"`
	Password string `env:"PGPASS" env-default:"password"`
	SSLMode  bool   `env:"PGSSLENABLED" env-default:"true"`
}

type AppConfig struct {
	GCloudProjectID      string `env:"GCLOUD_PROJECTID"`
	GCloudSubscriptionID string `env:"GCLOUD_SUBSCRIPTIONID"`
}

func main() {
	var dbConfig DatabaseConfig
	err := cleanenv.ReadConfig("config.env", &dbConfig)
	if err != nil {
		log.Fatal("Failed to read DatabaseConfig from config.env")
	}

	var appConfig AppConfig
	err = cleanenv.ReadConfig("config.env", &appConfig)
	if err != nil {
		log.Fatal("Failed to read AppConfig from config.env")
	}

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	client, err := pubsub.NewClient(context.Background(), appConfig.GCloudProjectID)
	if err != nil {
		log.Fatal("Failed to create a new secretmanager client")
	}
	defer client.Close()

	sub := client.Subscription(appConfig.GCloudSubscriptionID)
	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Fatalf("Error on receiving message with err %v", err)
	}

	http.ListenAndServe(":3000", r)
}
