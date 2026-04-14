package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Collinsthegreat/hng14_stage0_backend/internal/client"
	"github.com/Collinsthegreat/hng14_stage0_backend/internal/handler"
	"github.com/Collinsthegreat/hng14_stage0_backend/internal/service"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := getEnv("PORT", "8080")
	genderizeBaseURL := getEnv("GENDERIZE_BASE_URL", "https://api.genderize.io")
	timeoutStr := getEnv("HTTP_TIMEOUT_SECONDS", "5")

	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil {
		logger.Error("invalid HTTP_TIMEOUT_SECONDS", "error", err)
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	genderizeClient := client.NewGenderizeClient(httpClient, genderizeBaseURL)
	classifyService := service.NewClassifyService(genderizeClient)
	classifyHandler := handler.NewClassifyHandler(classifyService)

	r := chi.NewRouter()
	handler.RegisterRoutes(r, classifyHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("server stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
