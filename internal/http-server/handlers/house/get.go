package house

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
)

// TODO: разнести модели на слои: репо - entity, сервис - model и хэндлеры - DTO
type FlatsGetter interface {
	GetFlatsByHouseID(ctx context.Context, houseID int64) ([]models.Flat, error)
	GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]models.Flat, error)
}

type Request struct {
	ID int64 `json:"id" validate:"required,min=1"`
}

type Response struct {
	Flats []models.Flat `json:"flats" validate:"required,dive"`
}

func New(log *slog.Logger, flatsGetter FlatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.house.get.new"

		userType, ok := r.Context().Value("user_type").(string)
		if !ok {
			log.Error("user_type not found in context", slog.String("op", op))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var houseIDStr = chi.URLParam(r, "id")
		houseID, err := strconv.ParseInt(houseIDStr, 10, 64)
		if err != nil {
			log.Error("invalid house ID", slog.String("house_id", houseIDStr), slog.String("op", op))
			http.Error(w, "Invalid house ID", http.StatusBadRequest)
			return
		}

		var flats []models.Flat
		var response Response

		if userType == "moderator" {
			flats, err = flatsGetter.GetFlatsByHouseID(r.Context(), houseID)
		} else if userType == "client" {
			flats, err = flatsGetter.GetApprovedFlatsByHouseID(r.Context(), houseID)
		} else {
			log.Error("unauthorized access attempt", slog.String("user_type", userType), slog.String("op", op))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Error("failed to get flats", slog.String("op", op), sl.Err(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   err.Error(),
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}
			render.JSON(w, r, response)
			return
		}

		if len(flats) == 0 {
			flats = []models.Flat{}
		}

		response = Response{
			Flats: flats,
		}

		render.JSON(w, r, response)
		log.Info("request handled successfully", slog.String("op", op), slog.String("user_type", userType))
	}
}
