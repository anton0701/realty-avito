package login

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"

	"realty-avito/internal/converter"
	"realty-avito/internal/http-server/handlers"
	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories/usersRepo"
)

var jwtKey = []byte("jwt_most_secret_key")

type Claims struct {
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func LoginHandler(log *slog.Logger, usersRepository usersRepo.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.LoginHandler"

		ctx := r.Context()

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		var req handlers.LoginRequest

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

		log.Info("request pass", slog.Any("request", req))

		creds, err := converter.ConvertLoginRequestToUserCredentials(req)
		if err != nil {
			log.Error("failed to make pass hash", slog.String("op", op), slog.StringValue(err.Error()))

			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   "failed to login",
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}

			render.JSON(w, r, response)
			return
		}

		user, err := usersRepository.GetUserByCredentials(ctx, *creds)
		if err != nil {
			if errors.Is(err, usersRepo.ErrUserNotFound) {
				http.Error(w, "Пользователь не найден", http.StatusNotFound)
				return
			} else {
				log.Error("failed to login", slog.String("op", op), slog.StringValue(err.Error()))

				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusInternalServerError)

				response := models.InternalServerErrorResponse{
					Message:   "failed to login",
					RequestID: middleware.GetReqID(r.Context()),
					Code:      12345,
				}

				render.JSON(w, r, response)
				return
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			log.Error("failed to login", slog.String("op", op), slog.StringValue(err.Error()))

			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   "failed to login",
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}

			render.JSON(w, r, response)
			return
		}

		token, err := generateJWT(user.UserType, user.ID)
		if err != nil {
			log.Error("failed to create token", slog.String("op", op), slog.StringValue(err.Error()))

			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   err.Error(),
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}

			render.JSON(w, r, response)
			return
		}

		loginResponse := handlers.LoginResponse{
			Token: token,
		}

		if err := validator.New().Struct(loginResponse); err != nil {
			log.Error("failed to validate response body", sl.Err(err))
			http.Error(w, "failed to validate response body", http.StatusBadRequest)
			return
		}

		render.JSON(w, r, loginResponse)
	}
}

func generateJWT(userType string, userID int64) (string, error) {
	claims := &Claims{
		UserType:         userType,
		RegisteredClaims: jwt.RegisteredClaims{ID: fmt.Sprintf("%d", userID)},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
