package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
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

	projectID := "global-walker-125521"
	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatal("Failed to create a new secretmanager client")
	}
	defer client.Close()

	sub := client.Subscription("nest-monitor-subscription")
	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Fatalf("Error on receiving message with err %v", err)
	}

	http.ListenAndServe(":3000", r)
}
