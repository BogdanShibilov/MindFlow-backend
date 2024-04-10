package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/config"
	router "github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/auth"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/services/expert"
	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/storage/postgres"
	"github.com/bogdanshibilov/mindflowbackend/internal/pkg/httpserver"

	"github.com/gin-gonic/gin"
)

type App struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Run() {
	const op = "app.Run"

	pgDb, err := postgres.New(os.Getenv("PG_CONN_URL"))
	if err != nil {
		panic(op + " " + err.Error())
	}
	defer func() {
		err := pgDb.Close()
		if err != nil {
			a.log.Error(op, err)
		}
	}()

	auth := auth.New(
		pgDb,
		pgDb,
		a.cfg.TokenTTL,
		a.cfg.Secret,
	)

	experts := expert.New(pgDb, pgDb)

	handler := gin.New()
	router.New(handler, a.log, auth, experts)
	httpserver := httpserver.New(
		handler,
		httpserver.Port(a.cfg.Port),
	)
	httpserver.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		a.log.Info(op + " " + s.String())
	case err := <-httpserver.Notify():
		a.log.Error("%s: %w", op, err)
	}

	if err := httpserver.Shutdown(); err != nil {
		a.log.Error("%s: %w", op, err)
	}
}
