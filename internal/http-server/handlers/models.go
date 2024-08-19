package handlers

import "time"

type CreateFlatRequest struct {
	HouseID int64 `json:"house_id" validate:"required,min=1"`
	Price   int64 `json:"price" validate:"required,min=0"`
	Rooms   int64 `json:"rooms" validate:"required,min=1"`
}

type CreateFlatResponse struct {
	ID      int64                `json:"id" validate:"required,min=1"`
	HouseID int64                `json:"house_id" validate:"required,min=1"`
	Price   int64                `json:"price" validate:"required,min=0"`
	Rooms   int64                `json:"rooms" validate:"required,min=1"`
	Status  FlatModerationStatus `json:"status" validate:"required,oneof='created' 'approved' 'declined' 'on moderation'"`
}

type UpdateFlatRequest struct {
	ID     int64                `json:"id" validate:"required,min=1"`
	Status FlatModerationStatus `json:"status" validate:"required,oneof='created' 'approved' 'declined' 'on moderation'"`
}

type UpdateFlatResponse struct {
	ID      int64                `json:"id"`
	HouseID int64                `json:"house_id"`
	Price   int64                `json:"price"`
	Rooms   int64                `json:"rooms"`
	Status  FlatModerationStatus `json:"status"`
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

type House struct {
	ID        int64      `json:"id"`
	Address   string     `json:"address"`
	Year      int        `json:"year"`
	Developer *string    `json:"developer"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateHouseRequest struct {
	Address   string  `json:"address" validate:"required,min=1"`
	Year      int     `json:"year" validate:"required,min=1"`
	Developer *string `json:"developer,omitempty"`
}

type CreateHouseResponse struct {
	ID        int64   `json:"id"`
	Address   string  `json:"address"`
	Year      int     `json:"year"`
	Developer *string `json:"developer,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	UserType  string    `json:"user_type" validate:"required,oneof=client moderator"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	UserType string `json:"user_type" validate:"required,oneof=client moderator"`
}

type LoginRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
