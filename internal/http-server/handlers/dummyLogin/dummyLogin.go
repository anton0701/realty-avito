package dummyLogin

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	myMiddleware "realty-avito/internal/http-server/middleware"
	"realty-avito/internal/models"
)

type DummyLoginResponse struct {
	Token string `json:"token"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.dummyLogin.dummyLogin"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userType := r.URL.Query().Get("user_type")

		if userType != "client" && userType != "moderator" {
			log.Info("invalid user_type:", slog.StringValue(userType))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := myMiddleware.GenerateDummyJWT(userType)
		if err != nil {
			log.Error("failed to create dummy token", slog.String("op", op), slog.StringValue(err.Error()))

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

		render.JSON(w, r, DummyLoginResponse{
			Token: token,
		})
	}
}
