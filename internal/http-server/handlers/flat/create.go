package flat

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"io"
	"net/http"

	"realty-avito/internal/lib/logger/sl"
	"realty-avito/internal/models"
)

type FlatsWriter interface {
	CreateFlat(ctx context.Context, flatModel models.CreateFlatEntity) (*models.Flat, error)
}

type Request struct {
	HouseID int64 `json:"house_id" validate:"required,min=1"`
	Price   int64 `json:"price" validate:"required,min=0"`
	Rooms   int64 `json:"rooms" validate:"required,min=1"`
}

type Response struct {
	models.Flat `validate:"required,dive"`
}

func CreateFlatHandler(log *slog.Logger, flatsWriter FlatsWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.create.new"

		var req Request

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

		var response Response

		// TODO: перенести в конвертер
		flatEntity := models.CreateFlatEntity{
			HouseID: req.HouseID,
			Price:   req.Price,
			Rooms:   req.Rooms,
		}

		createdFlatModel, err := flatsWriter.CreateFlat(r.Context(), flatEntity)

		if err != nil {
			log.Error("failed to create flat", slog.String("op", op), sl.Err(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   err.Error(),
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}
			render.JSON(w, r, response)
			return
		}

		response = Response{*createdFlatModel}

		render.JSON(w, r, response)
		log.Info("request handled successfully", slog.String("op", op))
	}
}
