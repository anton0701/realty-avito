package models

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
