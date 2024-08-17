package flat

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"realty-avito/internal/converter"
	"realty-avito/internal/http-server/handlers"
	"realty-avito/internal/models"
	"realty-avito/internal/repositories/flat"
)

func UpdateFlatHandler(log *slog.Logger, flatsRepository flat.FlatsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.update"

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

		moderatorIDFromRequest, ok := r.Context().Value("moderator_id").(string)
		if !ok {
			log.Error("Error: no moderator_id in context", slog.String("op", op))
			w.WriteHeader(http.StatusInternalServerError)

			response := models.InternalServerErrorResponse{
				Message:   "Error: no moderator_id in context",
				RequestID: middleware.GetReqID(r.Context()),
				Code:      12345,
			}

			render.JSON(w, r, response)
			return
		}

		// TODO: вынести в func checkModerator()
		flatToUpdate, err := flatsRepository.GetFlatByFlatID(r.Context(), req.ID)
		if err != nil {
			// TODO: создать 1 метод, который это все заполняет
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

		if flatToUpdate.ModeratorID != nil && *flatToUpdate.ModeratorID != moderatorIDFromRequest {
			log.Error("failed to update flat, flat is under moderation by another moderator",
				slog.String("op", op))
			http.Error(w,
				"failed to update flat, flat is under moderation by another moderator",
				http.StatusForbidden)
			return
		}

		entityToUpdate := converter.ConvertUpdateFlatRequestToEntity(req)
		entityToUpdate.ModeratorID = &moderatorIDFromRequest
		now := time.Now()
		entityToUpdate.UpdatedAt = &now

		updatedFlat, err := flatsRepository.UpdateFlat(r.Context(), entityToUpdate)
		if err != nil {
			log.Error("failed to update flat",
				slog.String("op", op),
				slog.StringValue(err.Error()))

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
