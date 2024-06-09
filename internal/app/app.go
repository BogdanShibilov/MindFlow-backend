package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/mindflowbackend/internal/config"
	v1 "github.com/bogdanshibilov/mindflowbackend/internal/controller/http/v1"
	"github.com/bogdanshibilov/mindflowbackend/internal/db/postgres"
	"github.com/bogdanshibilov/mindflowbackend/internal/httpserver"
	"github.com/bogdanshibilov/mindflowbackend/internal/repository"
	authservice "github.com/bogdanshibilov/mindflowbackend/internal/services/auth"
	consultationservice "github.com/bogdanshibilov/mindflowbackend/internal/services/consultation"
	expertservice "github.com/bogdanshibilov/mindflowbackend/internal/services/expert"
	userservice "github.com/bogdanshibilov/mindflowbackend/internal/services/user"
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

	db, err := postgres.New(os.Getenv("PG_CONN_URL"))
	if err != nil {
		panic(op + " " + err.Error())
	}
	defer db.Close()

	userRepo := repository.NewUser(db)
	users := userservice.New(userRepo)
	auth := authservice.New(users, os.Getenv("JWTSECRET"), a.cfg.TokenTTL)
	expertsRepo := repository.NewExpert(db)
	experts := expertservice.New(expertsRepo, userRepo)
	consultRepo := repository.NewConsultation(db)
	consultations := consultationservice.New(*consultRepo, *userRepo)

	handler := gin.New()
	v1.NewRouter(handler, a.log, auth, experts, users, consultations)
	httpserver := httpserver.New(handler, httpserver.Port(a.cfg.Port))
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
