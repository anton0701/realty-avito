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

// TODO: вроде правильно
func New(log *slog.Logger) http.HandlerFunc {
	// TODO: вынести строки в константы
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.dummyLogin.dummyLogin.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userType := r.URL.Query().Get("user_type")
		// TODO: сделать через попытку кастануть к UserType (если не получается, то error)
		if userType != "client" && userType != "moderator" {
			log.Info("invalid user_type:", slog.StringValue(userType))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := myMiddleware.GenerateJWT(userType)
		if err != nil {
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   "Что-то пошло не так",
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
