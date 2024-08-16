package house

import "time"

type CreateHouseEntity struct {
	Address   string  `json:"address" validate:"required"`
	Year      int     `json:"year" validate:"required"`
	Developer *string `json:"developer"`
}

type HouseEntity struct {
	ID        int64
	Address   string
	Year      int
	Developer *string
	CreatedAt time.Time
}
