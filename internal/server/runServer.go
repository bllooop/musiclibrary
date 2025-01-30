package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/bllooop/musiclibrary/internal/delivery/api"
	"github.com/bllooop/musiclibrary/internal/repository"
	"github.com/bllooop/musiclibrary/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func Run() {
	logger := zerolog.New(os.Stdout).Level(zerolog.TraceLevel)
	//errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with env")
	}
	dbpool, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	})
	migratePath := "./migrations"
	if err := repository.RunMigrate(repository.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	}, migratePath); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error when migrating")
	}
	if err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error with database")
	}

	repos := repository.NewRepository(dbpool)
	usecases := usecase.NewService(repos)
	handler := handlers.NewHandler(usecases)
	srv := new(Server)

	go func() {
		if err := srv.RunServer(os.Getenv("SERVERPORT"), handler.InitRoutes()); err != nil && err == http.ErrServerClosed {
			logger.Info().Msg("Server was shut down gracefully")
		} else {
			logger.Error().Err(err).Msg("")
			logger.Fatal().Msg("There was an error when starting the server")
		}
	}()

	logger.Info().Msg("Server is running")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logger.Info().Msg("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer dbpool.Close()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("")
		logger.Fatal().Msg("There was an error while shutting down the server")
	}
}
