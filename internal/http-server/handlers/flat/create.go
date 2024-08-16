package flat

import (
	"context"
	"errors"
	"io"
	"net/http"
	"realty-avito/internal/http-server/handlers"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"realty-avito/internal/converter"
	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories/flat"
)

type FlatsWriter interface {
	CreateFlat(ctx context.Context, flatModel flat.CreateFlatEntity) (*flat.FlatEntity, error)
	UpdateFlat(ctx context.Context, updateFlatModel flat.UpdateFlatEntity) (*flat.FlatEntity, error)
}

func CreateFlatHandler(log *slog.Logger, flatsWriter FlatsWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.create"

		var req handlers.CreateFlatRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", slog.String("request_id", middleware.GetReqID(r.Context())), slog.String("op", op))
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
			//validateErr := err.(validator.ValidationErrors)
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		flatEntity := converter.ConvertCreateFlatRequestToEntity(req)

		createdFlatEntity, err := flatsWriter.CreateFlat(r.Context(), flatEntity)
		if err != nil {
			log.Error("failed to create flat", slog.String("op", op), sl.Err(err))

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

		response := converter.ConvertFlatEntityToCreateResponse(createdFlatEntity)

		render.JSON(w, r, response)
		log.Info("request handled successfully", slog.String("op", op))
	}
}
