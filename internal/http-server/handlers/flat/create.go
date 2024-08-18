package flat

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"realty-avito/internal/client/db"
	"realty-avito/internal/converter"
	"realty-avito/internal/http-server/handlers"
	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories/flatsRepo"
)

type FlatsWriter interface {
	CreateFlat(ctx context.Context, flatModel flatsRepo.CreateFlatEntity) (*flatsRepo.FlatEntity, error)
	UpdateFlat(ctx context.Context, updateFlatModel flatsRepo.UpdateFlatEntity) (*flatsRepo.FlatEntity, error)
}

type HousesWriter interface {
	UpdateHouseUpdatedAt(ctx context.Context, houseID int64) error
}

func CreateFlatHandler(log *slog.Logger, flatsWriter FlatsWriter, housesWriter HousesWriter, txManager db.TxManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.create"

		ctx := r.Context()

		var req handlers.CreateFlatRequest

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
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		flatEntity := converter.ConvertCreateFlatRequestToEntity(req)

		var successfullyCreatedFlatEntity *flatsRepo.FlatEntity

		err = txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			var errTx error
			createdFlatEntity, errTx := flatsWriter.CreateFlat(ctx, flatEntity)
			if errTx != nil {
				return errTx
			}

			errTx = housesWriter.UpdateHouseUpdatedAt(ctx, req.HouseID)
			if errTx != nil {
				return errTx
			}

			successfullyCreatedFlatEntity = createdFlatEntity

			return nil
		})

		if err != nil {
			log.Error("failed to create flat", slog.String("op", op), sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Retry-After", "60")

			response := models.InternalServerErrorResponse{
				Message:   err.Error(),
				RequestID: middleware.GetReqID(ctx),
				Code:      12345,
			}
			render.JSON(w, r, response)
			return
		}

		response := converter.ConvertFlatEntityToCreateResponse(successfullyCreatedFlatEntity)

		render.JSON(w, r, response)
		log.Info("request handled successfully", slog.String("op", op))
	}
}
