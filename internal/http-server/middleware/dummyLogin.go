package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"realty-avito/internal/models"
)

type DummyLoginResponse struct {
	Token string `json:"token"`
}

func DummyLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: вынести строки в константы
	userType := r.URL.Query().Get("user_type")

	if userType != "client" && userType != "moderator" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := GenerateJWT(userType)
	if err != nil {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusInternalServerError)

		response := models.InternalServerErrorResponse{
			Message:   "что-то пошло не так",
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
