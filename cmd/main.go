package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slog"

	"realty-avito/internal/config"
	"realty-avito/internal/http-server/handlers/dummyLogin"
	"realty-avito/internal/http-server/handlers/house"
	myMiddleware "realty-avito/internal/http-server/middleware"
	mwLogger "realty-avito/internal/http-server/middleware/logger"
	"realty-avito/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	log.Info(
		"starting realty-avito",
		slog.String("env", cfg.Env),
	)

	// TODO: init storage postgres
	pool, err := initPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Error("failed to initialize postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// TODO: создать flatsRepo + housesRepo + возможно moderatedFlatsRepo

	// TODO: init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//router.Get("/dummyLogin", dummyLogin.DummyLogin)
	router.Get("/dummyLogin", dummyLogin.New(log))
	// TODO: продолжить завтра
	//router.Get("/house/{id}", house.New(log, flatsRepository))
	//router.Use(myMiddleware.JWTMiddleware())

	router.With(myMiddleware.JWTMiddleware()).Get("/house/{id}", func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value("user_type").(string)

		log.Info("user_type", slog.String(userType, userType))

		if userType == "moderator" {
			house.GetModeratorHouseHandler()(w, r)
		} else if userType == "client" {
			house.GetClientHouseHandler()(w, r)
		}
	})

	// TODO: run server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("server is listening", slog.String("address", cfg.HTTPServer.Address))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("could not listen on", slog.String("address", cfg.HTTPServer.Address), slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func initPostgres(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PostgreSQL DSN: %w", err)
	}

	dbpool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to PostgreSQL: %w", err)
	}

	return dbpool, nil
}
