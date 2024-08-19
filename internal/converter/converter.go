package converter

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	handlers "realty-avito/internal/http-server/handlers"
	flatRepo "realty-avito/internal/repositories/flatsRepo"
	houseRepo "realty-avito/internal/repositories/housesRepo"
	"realty-avito/internal/repositories/usersRepo"
)

func ConvertCreateFlatRequestToEntity(req handlers.CreateFlatRequest) flatRepo.CreateFlatEntity {
	return flatRepo.CreateFlatEntity{
		HouseID: req.HouseID,
		Price:   req.Price,
		Rooms:   req.Rooms,
		Status:  flatRepo.StatusCreated,
	}
}

func ConvertUpdateFlatRequestToEntity(req handlers.UpdateFlatRequest) flatRepo.UpdateFlatEntity {
	return flatRepo.UpdateFlatEntity{
		ID:     req.ID,
		Status: flatRepo.FlatModerationStatus(req.Status),
	}
}

func ConvertFlatEntityToCreateResponse(entity *flatRepo.FlatEntity) handlers.CreateFlatResponse {
	return handlers.CreateFlatResponse{
		ID:      entity.ID,
		HouseID: entity.HouseID,
		Price:   entity.Price,
		Rooms:   entity.Rooms,
		Status:  handlers.FlatModerationStatus(entity.Status),
	}
}

func ConvertFlatEntityToUpdateResponse(entity *flatRepo.FlatEntity) handlers.UpdateFlatResponse {
	return handlers.UpdateFlatResponse{
		ID:      entity.ID,
		HouseID: entity.HouseID,
		Price:   entity.Price,
		Rooms:   entity.Rooms,
		Status:  handlers.FlatModerationStatus(entity.Status),
	}
}

func ConvertCreateHouseRequestToEntity(req handlers.CreateHouseRequest) houseRepo.CreateHouseEntity {
	return houseRepo.CreateHouseEntity{
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
	}
}

func ConvertEntityToCreateHouseResponse(entity *houseRepo.HouseEntity) handlers.CreateHouseResponse {
	return handlers.CreateHouseResponse{
		ID:        entity.ID,
		Address:   entity.Address,
		Year:      entity.Year,
		Developer: entity.Developer,
		CreatedAt: entity.CreatedAt.Format(time.RFC3339),
	}
}

func ConvertFlatEntitiesToFlats(entities []flatRepo.FlatEntity) []handlers.Flat {
	flats := make([]handlers.Flat, len(entities))

	for i, entity := range entities {
		flats[i] = ConvertEntityToFlat(entity)
	}
	return flats
}

func ConvertEntityToFlat(entity flatRepo.FlatEntity) handlers.Flat {
	return handlers.Flat{
		ID:      entity.ID,
		HouseID: entity.HouseID,
		Price:   entity.Price,
		Rooms:   entity.Rooms,
		Status:  handlers.FlatModerationStatus(entity.Status),
	}
}

func ConvertRegisterRequestToUserEntity(req handlers.RegisterRequest) (*usersRepo.UserEntity, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &usersRepo.UserEntity{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		UserType:     req.UserType,
		CreatedAt:    time.Now(),
	}, nil
}

func ConvertLoginRequestToUserCredentials(req handlers.LoginRequest) (*usersRepo.UserCredentials, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	return &usersRepo.UserCredentials{
		ID:           req.ID,
		PasswordHash: hashedPassword,
	}, nil
}

func ConvertUserToUserEntity(user handlers.User) (*usersRepo.UserEntity, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	return &usersRepo.UserEntity{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		UserType:     user.UserType,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
