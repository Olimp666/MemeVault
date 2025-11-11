package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Olimp666/MemeVault/internal/api"
	"github.com/Olimp666/MemeVault/internal/repository"
	"github.com/Olimp666/MemeVault/internal/service"
)

type Config struct {
	Database struct {
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		Name     string `env:"POSTGRES_DB"`
		Host     string `env:"POSTGRES_HOST"`
		Port     int64  `env:"POSTGRES_PORT"`
	}

	ServerURL string `env:"SERVER_URL"`
}

type Application struct {
	cfg    *Config
	db     *sqlx.DB
	server *http.Server

	errChan chan error

	wg sync.WaitGroup
}

func New() *Application {
	return &Application{
		errChan: make(chan error),
	}
}

func (a *Application) Start(_ context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}

	if err := a.initDB(); err != nil {
		return fmt.Errorf("can't init database connection: %w", err)
	}

	if err := a.initServer(); err != nil {
		return fmt.Errorf("can't init server: %w", err)
	}

	return nil
}

func (a *Application) Wait(ctx context.Context, cancel context.CancelFunc) error {
	var appErr error

	errWg := sync.WaitGroup{}

	errWg.Add(1)

	go func() {
		defer errWg.Done()

		for err := range a.errChan {
			cancel()

			fmt.Println(err)
			appErr = err
		}
	}()

	<-ctx.Done()
	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}

func (a *Application) initConfig() error {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("can't parse env vars: %w", err)
	}

	a.cfg = cfg

	return nil
}

func (a *Application) initDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		a.cfg.Database.Host,
		a.cfg.Database.Port,
		a.cfg.Database.User,
		a.cfg.Database.Password,
		a.cfg.Database.Name)

	dbConn, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	if err = dbConn.Ping(); err != nil {
		return err
	}

	a.db = dbConn
	log.Println("Database connection established")

	return nil
}

func (a *Application) initServer() error {
	repo := repository.NewRepository(a.db)
	svc := service.NewService(repo)
	handler := api.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handler.UploadImage)
	mux.HandleFunc("/images", handler.ImagesByTags)
	mux.HandleFunc("/user/images", handler.ImagesByUser)

	a.server = &http.Server{
		Addr:    a.cfg.ServerURL,
		Handler: mux,
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		log.Printf("Starting HTTP server on %s", a.cfg.ServerURL)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	return nil
}
