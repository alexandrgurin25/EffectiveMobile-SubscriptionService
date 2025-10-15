package main

import (
	"context"
	"net/http"
	"subscriptions/internal/config"
	"subscriptions/internal/repositories"
	service "subscriptions/internal/services"
	"subscriptions/internal/transport/http/handlers"
	"subscriptions/pkg/logger"
	"subscriptions/pkg/postgres"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(ctx, "unable to load config", zap.Error(err))
		return
	}

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "unable to connect db", zap.Error(err))
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Successful start!")

	r := chi.NewRouter()
	repository := repositories.New(db)

	service := service.New(repository)

	handlers := handlers.New(service)

	r.Route("/api/subscriptions/", func(r chi.Router) {
		r.Post("/", handlers.Create)
		r.Get("/", handlers.GetList) // /api/subscriptions?page=1&limit=10
		r.Get("/{id}", handlers.Get)
		r.Put("/{id}", handlers.Put)
		r.Delete("/{id}", handlers.Delete)
		r.Get("/summary/{user_id}/{service_name}", handlers.GetSummary)
	})

	http.ListenAndServe(":8080", r)
}
