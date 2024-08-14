package house

import (
	"github.com/go-chi/render"
	"net/http"

	"golang.org/x/exp/slog"
)

// TODO: это будет FlatsRepository
type FlatsGetter interface {
	GetFlats(houseID int64) []Flat
}

type Request struct {
	ID int64 `json:"id" validate:"required,min=1"`
}

type Response struct {
	Flats []Flat `json:"flats" validate:"required,dive"`
}

type FlatModerationStatus string

const (
	StatusCreated      FlatModerationStatus = "created"
	StatusApproved     FlatModerationStatus = "approved"
	StatusDeclined     FlatModerationStatus = "declined"
	StatusOnModeration FlatModerationStatus = "on moderation"
)

type Flat struct {
	ID      int64                `json:"id" validate:"required,min=1"`
	HouseID int64                `json:"house_id" validate:"required,min=1"`
	Price   int64                `json:"price" validate:"required,min=0"`
	Rooms   int64                `json:"rooms" validate:"required,min=1"`
	Status  FlatModerationStatus `json:"status" validate:"required,oneof='created approved declined \"on moderation\"'"`
}

func New(log *slog.Logger /*, flatsGetter FlatsGetter*/) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.house.get.new"

		userType, ok := r.Context().Value("user_type").(string)
		if !ok {
			log.Error("user_type not found in context", slog.String("op", op))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var response Response

		if userType == "moderator" {
			// TODO: implement me!
			flats := []Flat{
				{ID: 1, HouseID: 101, Price: 500000, Rooms: 2, Status: StatusApproved},
				{ID: 2, HouseID: 101, Price: 600000, Rooms: 3, Status: StatusApproved},
				{ID: 3, HouseID: 101, Price: 700000, Rooms: 4, Status: StatusOnModeration},
				{ID: 4, HouseID: 101, Price: 800000, Rooms: 5, Status: StatusDeclined},
				{ID: 5, HouseID: 101, Price: 900000, Rooms: 6, Status: StatusCreated},
			}

			response = Response{
				Flats: flats,
			}

		} else if userType == "client" {
			// TODO: implement me!
			flat := Flat{
				ID:      1,
				HouseID: 101,
				Price:   500000,
				Rooms:   2,
				Status:  StatusApproved,
			}

			response = Response{
				Flats: []Flat{flat},
			}
		} else {
			log.Error("unauthorized access attempt", slog.String("user_type", userType), slog.String("op", op))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		render.JSON(w, r, response)
		log.Info("request handled successfully", slog.String("op", op), slog.String("user_type", userType))
	}
}
