package models

import "time"

type UserType string

const (
	Client    UserType = "client"
	Moderator UserType = "moderator"
)

type InternalServerErrorResponse struct {
	Message   string `json:"message" validate:"required"`
	RequestID string `json:"request_id,omitempty"`
	Code      int    `json:"code,omitempty"`
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
	Status  FlatModerationStatus `json:"status" validate:"required,oneof='created' 'approved' 'declined' 'on moderation'"`
}

type CreateFlatEntity struct {
	HouseID int64 `json:"house_id" validate:"required,min=1"`
	Price   int64 `json:"price" validate:"required,min=0"`
	Rooms   int64 `json:"rooms" validate:"required,min=1"`
}

type UpdateFlatEntity struct {
	ID     int64                `json:"id" validate:"required,min=1"`
	Status FlatModerationStatus `json:"status" validate:"required,oneof='created' 'approved' 'declined' 'on moderation'"`
}

type House struct {
	ID        int64      `json:"id"`
	Address   string     `json:"address"`
	Year      int        `json:"year"`
	Developer *string    `json:"developer"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateHouseEntity struct {
	Address   string  `json:"address" validate:"required"`
	Year      int     `json:"year" validate:"required"`
	Developer *string `json:"developer"`
}
