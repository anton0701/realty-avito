package main

import (
	"os"

	"golang.org/x/exp/slog"

	"realty-avito/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	log.Info(
		"starting realty-avito",
		slog.String("env", cfg.Env),
	)

	// TODO: init storage

	// TODO: init router

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		//log = setupPrettySlog()
		log = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
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
