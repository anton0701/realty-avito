package register

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"realty-avito/internal/converter"
	"realty-avito/internal/http-server/handlers"
	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories/usersRepo"
)

func RegisterHandler(log *slog.Logger, usersRepository usersRepo.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.RegisterHandler"

		ctx := r.Context()

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		var req handlers.RegisterRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty",
				slog.String("request_id", middleware.GetReqID(ctx)),
				slog.String("op", op))
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request body", sl.Err(err))
			http.Error(w, "failed to validate request body", http.StatusBadRequest)
			return
		}

		userEntity, err := converter.ConvertRegisterRequestToUserEntity(req)
		if err != nil {
			log.Error("failed to register user", slog.String("op", op), slog.StringValue(err.Error()))

			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   "failed to register user",
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}

			render.JSON(w, r, response)
			return
		}

		createdUser, err := usersRepository.CreateUser(ctx, *userEntity)
		if err != nil {
			if errors.Is(err, usersRepo.ErrEmailExists) {
				log.Error("failed to create user: email already exist", sl.Err(err))
				http.Error(w, "failed to validate request body", http.StatusBadRequest)
				return
			}
		}

		// Подготовка ответа
		response := handlers.RegisterResponse{
			UserID: createdUser.UUID,
		}

		render.JSON(w, r, response)
	}
}
