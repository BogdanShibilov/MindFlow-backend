package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bogdanshibilov/mindflowbackend/internal/mindflow/config"
	router "github.com/bogdanshibilov/mindflowbackend/internal/mindflow/controller"
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

	handler := gin.New()
	router.New(handler)
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
