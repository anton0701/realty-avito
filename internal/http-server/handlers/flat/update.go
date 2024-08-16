package flat

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	"realty-avito/internal/converter"
	"realty-avito/internal/http-server/handlers"

	"realty-avito/internal/models"
)

func UpdateFlatHandler(log *slog.Logger, flatsWriter FlatsWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.update"

		// TODO: разнести entity и request + converter + validator
		var req handlers.UpdateFlatRequest

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

		entityToUpdate := converter.ConvertUpdateFlatRequestToEntity(req)

		updatedFlat, err := flatsWriter.UpdateFlat(r.Context(), entityToUpdate)
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

		response := converter.ConvertFlatEntityToUpdateResponse(updatedFlat)

		render.JSON(w, r, response)
		log.Info("flat updated successfully", slog.Int64("flat_id", response.ID))
	}
}
