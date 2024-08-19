package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"

	"realty-avito/internal/client/db/pg"
	"realty-avito/internal/client/db/transaction"
	"realty-avito/internal/config"
	"realty-avito/internal/http-server/handlers/dummyLogin"
	"realty-avito/internal/http-server/handlers/flat"
	"realty-avito/internal/http-server/handlers/house"
	"realty-avito/internal/http-server/handlers/login"
	"realty-avito/internal/http-server/handlers/register"
	myMiddleware "realty-avito/internal/http-server/middleware"
	mwLogger "realty-avito/internal/http-server/middleware/logger"
	"realty-avito/internal/lib/logger"
	flatRepo "realty-avito/internal/repositories/flatsRepo"
	houseRepo "realty-avito/internal/repositories/housesRepo"
	"realty-avito/internal/repositories/usersRepo"
	"realty-avito/postgres"
)

func main() {
	ctx := context.Background()

	// init config
	cfg := config.MustLoad()

	// init logger: slog
	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"starting realty-avito",
		slog.String("env", cfg.Env),
	)

	dsn := postgres.CreatePostgresDSN(cfg.Postgres)
	pgClient, err := pg.New(ctx, dsn)
	if err != nil {
		log.Error("failed to initialize postgres client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer pgClient.Close()

	// init transaction manager
	txManager := transaction.NewTransactionManager(pgClient.DB())

	// init repo
	flatsRepo := flatRepo.NewFlatsRepository(pgClient)
	housesRepo := houseRepo.NewHousesRepository(pgClient)
	usersRepository := usersRepo.NewUserRepository(pgClient)

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// GET /dummyLogin
	router.Get("/dummyLogin", dummyLogin.New(log))

	// GET /house/{id}
	router.Route("/house/{id}", func(r chi.Router) {
		r.Use(myMiddleware.JWTMiddleware)
		r.Get("/", house.GetFlatsInHouseHandler(log, flatsRepo))
	})

	// POST /house/create
	router.Route("/house/create", func(r chi.Router) {
		r.Use(myMiddleware.JWTModeratorOnlyMiddleware)
		r.Post("/", house.CreateHouseHandler(log, housesRepo))
	})

	// POST /flatsRepo/create
	router.Route("/flat/create", func(r chi.Router) {
		r.Use(myMiddleware.JWTMiddleware)
		r.Post("/", flat.CreateFlatHandler(log, flatsRepo, housesRepo, txManager))
	})

	// POST /flatsRepo/update
	router.Route("/flat/update", func(r chi.Router) {
		r.Use(myMiddleware.JWTModeratorOnlyMiddleware)
		r.Post("/", flat.UpdateFlatHandler(log, flatsRepo))
	})

	// POST /register
	router.Post("/register", register.RegisterHandler(log, usersRepository))

	// POST /login
	router.Post("/login", login.LoginHandler(log, usersRepository))

	// Run server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("server is listening", slog.String("address", cfg.HTTPServer.Address))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("could not listen on",
			slog.String("address", cfg.HTTPServer.Address),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
