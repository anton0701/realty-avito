package house

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories"
)

type CreateHouseRequest struct {
	Address   string  `json:"address" validate:"required,min=1"`
	Year      int     `json:"year" validate:"required,min=1"`
	Developer *string `json:"developer,omitempty"`
}

type CreateHouseResponse struct {
	ID        int64   `json:"id"`
	Address   string  `json:"address"`
	Year      int     `json:"year"`
	Developer *string `json:"developer,omitempty"`
	CreatedAt string  `json:"created_at"`
}

func CreateHouseHandler(log *slog.Logger, housesRepo repositories.HousesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.house.create.new"

		var req CreateHouseRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error("failed to decode request body", slog.String("op", op), sl.Err(err))
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error("failed to decode request body", slog.String("op", op), slog.StringValue(err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation failed", slog.String("op", op), slog.StringValue(err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// TODO: через конвертер
		houseModel := models.CreateHouseEntity{
			Address:   req.Address,
			Year:      req.Year,
			Developer: req.Developer,
		}

		createdHouse, err := housesRepo.CreateHouse(r.Context(), houseModel)
		if err != nil {
			log.Error("failed to create house", slog.String("op", op), slog.StringValue(err.Error()))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			response := models.InternalServerErrorResponse{
				Message:   err.Error(),
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}
			render.JSON(w, r, response)
			return
		}

		response := CreateHouseResponse{
			ID:        createdHouse.ID,
			Address:   createdHouse.Address,
			Year:      createdHouse.Year,
			Developer: createdHouse.Developer,
			CreatedAt: createdHouse.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
		log.Info("house created successfully", slog.String("op", op), slog.Int64("house_id", response.ID))
	}
}
