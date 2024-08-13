package house

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
