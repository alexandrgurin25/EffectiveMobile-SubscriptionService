package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"subscriptions/internal/config"
	"subscriptions/internal/repositories"
	"subscriptions/internal/services"
	"subscriptions/internal/transport/http/handlers"
	"subscriptions/pkg/logger"
	"subscriptions/pkg/postgres"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "subscriptions/docs"
)

// @title Subscriptions Service API
// @version 1.0.0
// @description Сервис для управления подписками пользователей

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

	service := services.New(repository)

	handlers := handlers.New(service)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL для JSON документации
	))

	r.Route("/api/subscriptions/", func(r chi.Router) {
		r.Post("/", handlers.Create)
		r.Get("/", handlers.GetList) // /api/subscriptions?page=1&limit=10
		r.Get("/{id}", handlers.Get)
		r.Put("/{id}", handlers.Put)
		r.Delete("/{id}", handlers.Delete)
		r.Get("/summary/{user_id}/{service_name}", handlers.GetSummary)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second, 
		WriteTimeout: 15 * time.Second, 
		IdleTimeout:  60 * time.Second, 
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Error(ctx, "Server error", zap.Error(err))
		return

	case <-shutdown:
		log.Info(ctx, "Starting graceful shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error(ctx, "Graceful shutdown failed", zap.Error(err))
			server.Close()
		}
	}

	log.Info(ctx, "Server stopped")
}
