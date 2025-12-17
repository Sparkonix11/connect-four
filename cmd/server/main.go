package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"connect-four/internal/api"
	"connect-four/internal/database"
	"connect-four/internal/kafka"
	"connect-four/internal/models"
	"connect-four/pkg/config"
)

func main() {
	// Configure zerolog for pretty console output in development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load configuration
	cfg := config.Load()
	log.Info().Str("port", cfg.ServerPort).Msg("Starting Connect Four server")

	// Connect to database
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close(db)

	// Run migrations (auto-migrate GORM models)
	if err := models.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	log.Info().Msg("Database migrations completed")

	// Create Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopicEvents, cfg.KafkaEnabled, cfg.KafkaUsername, cfg.KafkaPassword)
	defer kafkaProducer.Close()

	// Create Kafka consumer for analytics
	kafkaConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopicEvents, "connect-four-analytics", cfg.KafkaEnabled, cfg.KafkaUsername, cfg.KafkaPassword)
	defer kafkaConsumer.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer in background
	go kafkaConsumer.Start(ctx, kafka.ProcessEvent)

	// Create API server with Kafka producer
	server := api.NewServer(db, cfg, kafkaProducer)
	server.Start()

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Configure for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:      c.Handler(server.Router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Msgf("Server listening on http://%s:%s", cfg.ServerHost, cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Cancel context to stop Kafka consumer
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
}
