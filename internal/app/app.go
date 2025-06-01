package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/romanpitatelev/denet/internal/configs"
	"github.com/romanpitatelev/denet/internal/controller/rest"
	taskshandler "github.com/romanpitatelev/denet/internal/controller/rest/tasks-handler"
	usershandler "github.com/romanpitatelev/denet/internal/controller/rest/users-handler"
	"github.com/romanpitatelev/denet/internal/repository/store"
	tasksrepo "github.com/romanpitatelev/denet/internal/repository/tasks-repo"
	usersrepo "github.com/romanpitatelev/denet/internal/repository/users-repo"
	tasksservice "github.com/romanpitatelev/denet/internal/usecase/task-service"
	usersservice "github.com/romanpitatelev/denet/internal/usecase/users-service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
)

func Run(cfg *configs.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	log.Level(level)

	db, err := store.New(ctx, store.Config{Dsn: cfg.PostgresDSN})
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to database")
	}

	if err := db.Migrate(migrate.Up); err != nil {
		log.Panic().Err(err).Msg("failed to migrate")
	}

	log.Info().Msg("successful migration")

	usersRepo := usersrepo.New(db)
	tasksRepo := tasksrepo.New(db)

	usersService := usersservice.New(usersRepo)
	tasksService := tasksservice.New(tasksRepo)

	usersHandler := usershandler.New(usersService)
	tasksHandler := taskshandler.New(tasksService)

	server := rest.New(
		rest.Config{BindAddress: cfg.BindAddress},
		usersHandler,
		tasksHandler,
		rest.GetPublicKey(),
	)

	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("failed to run the server: %w", err)
	}

	return nil
}
