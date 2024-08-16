package flat

type CreateFlatEntity struct {
	HouseID int64
	Price   int64
	Rooms   int64
	Status  FlatModerationStatus
}

type UpdateFlatEntity struct {
	ID          int64
	Status      FlatModerationStatus
	ModeratorID *string
}

type FlatEntity struct {
	ID          int64
	HouseID     int64
	Price       int64
	Rooms       int64
	Status      FlatModerationStatus
	ModeratorID *string
}

type FlatModerationStatus string

const (
	StatusCreated      FlatModerationStatus = "created"
	StatusApproved     FlatModerationStatus = "approved"
	StatusDeclined     FlatModerationStatus = "declined"
	StatusOnModeration FlatModerationStatus = "on moderation"
)
