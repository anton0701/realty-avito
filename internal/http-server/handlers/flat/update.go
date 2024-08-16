package flat

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"realty-avito/internal/models"
)

type UpdateFlatRequest struct {
	ID     int64                       `json:"id" validate:"required,min=1"`
	Status models.FlatModerationStatus `json:"status" validate:"required,oneof='created' 'approved' 'declined' 'on moderation'"`
}

type UpdateFlatResponse struct {
	ID      int64  `json:"id"`
	HouseID int64  `json:"house_id"`
	Price   int64  `json:"price"`
	Rooms   int64  `json:"rooms"`
	Status  string `json:"status"`
}

func UpdateFlatHandler(log *slog.Logger, flatsWriter FlatsWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.update"

		// TODO: разнести entity и request + converter + validator
		var req models.UpdateFlatEntity
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error(err.Error())
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		updatedFlat, err := flatsWriter.UpdateFlat(r.Context(), req)
		if err != nil {
			log.Error("failed to update flat", slog.String("op", op), slog.StringValue(err.Error()))

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

		response := UpdateFlatResponse{
			ID:      updatedFlat.ID,
			HouseID: updatedFlat.HouseID,
			Price:   updatedFlat.Price,
			Rooms:   updatedFlat.Rooms,
			Status:  string(updatedFlat.Status),
		}

		render.JSON(w, r, response)
		log.Info("flat updated successfully", slog.Int64("flat_id", response.ID))
	}
}
