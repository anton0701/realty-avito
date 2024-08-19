package usersRepo

import "time"

type UserType string

const (
	UserTypeClient    UserType = "client"
	UserTypeModerator UserType = "moderator"
)

type UserEntity struct {
	ID           int64
	Email        string
	PasswordHash string
	UUID         string
	UserType     string
	CreatedAt    time.Time
}

type UserCredentials struct {
	ID           string
	PasswordHash string
}
