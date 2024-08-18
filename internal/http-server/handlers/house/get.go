package house

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"realty-avito/internal/converter"
	"realty-avito/internal/http-server/handlers"
	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories/flatsRepo"
)

type FlatsGetter interface {
	GetFlatsByHouseID(ctx context.Context, houseID int64) ([]flatsRepo.FlatEntity, error)
	GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]flatsRepo.FlatEntity, error)
}

type Request struct {
	ID int64 `json:"id" validate:"required,min=1"`
}

type Response struct {
	Flats []handlers.Flat `json:"flats" validate:"required,dive"`
}

func GetFlatsInHouseHandler(log *slog.Logger, flatsRepository flatsRepo.FlatsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.house.get"

		userType, ok := r.Context().Value("user_type").(string)
		if !ok {
			log.Error(
				"user_type not found in context",
				slog.String("op", op))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var houseIDStr = chi.URLParam(r, "id")
		houseID, err := strconv.ParseInt(houseIDStr, 10, 64)
		if err != nil {
			log.Error("invalid house ID",
				slog.String("house_id", houseIDStr),
				slog.String("op", op))
			http.Error(w, "Invalid house ID", http.StatusBadRequest)
			return
		}

		var flatEntities []flatsRepo.FlatEntity
		var response Response

		if userType == "moderator" {
			flatEntities, err = flatsRepository.GetFlatsByHouseID(r.Context(), houseID)
		} else if userType == "client" {
			flatEntities, err = flatsRepository.GetApprovedFlatsByHouseID(r.Context(), houseID)
		} else {
			log.Error("unauthorized access attempt",
				slog.String("user_type", userType),
				slog.String("op", op))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Error(
				"failed to get flats",
				slog.String("op", op),
				sl.Err(err),
			)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Retry-After", "60")

			response := models.InternalServerErrorResponse{
				Message:   err.Error(),
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}
			render.JSON(w, r, response)
			return
		}

		flats := converter.ConvertFlatEntitiesToFlats(flatEntities)

		if len(flats) == 0 {
			flats = []handlers.Flat{}
		}

		response = Response{
			Flats: flats,
		}

		render.JSON(w, r, response)
		log.Info(
			"request handled successfully",
			slog.String("op", op),
			slog.String("user_type", userType),
		)
	}
}
