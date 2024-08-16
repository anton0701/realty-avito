package converter

import (
	"time"

	handlers "realty-avito/internal/http-server/handlers"
	flatRepo "realty-avito/internal/repositories/flat"
	houseRepo "realty-avito/internal/repositories/house"
)

// TODO: возвращать не только конвертированную модель, но и ошибку !!!!
// TODO: возвращать ошибку, если не получилось конвертировать status !!!!
// TODO: в сервисе проверять ошибку на нил - возвращать 500 ошибку, если пришел эррор!!!

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
