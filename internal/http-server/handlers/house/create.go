package house

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
	"realty-avito/internal/repositories/house"
)

func CreateHouseHandler(log *slog.Logger, housesRepo house.HousesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.house.create"

		var req handlers.CreateHouseRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty",
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error("failed to decode request body",
				slog.String("op", op),
				sl.Err(err),
			)
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error("failed to decode request body",
				slog.String("op", op),
				slog.StringValue(err.Error()),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation failed",
				slog.String("op", op),
				slog.StringValue(err.Error()),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		houseEntity := converter.ConvertCreateHouseRequestToEntity(req)

		createdHouseEntity, err := housesRepo.CreateHouse(r.Context(), houseEntity)
		if err != nil {
			log.Error("failed to create house",
				slog.String("op", op),
				slog.StringValue(err.Error()),
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

		response := converter.ConvertEntityToCreateHouseResponse(createdHouseEntity)

		render.JSON(w, r, response)
		log.Info("house created successfully",
			slog.String("op", op),
			slog.Int64("house_id", response.ID),
		)
	}
}
