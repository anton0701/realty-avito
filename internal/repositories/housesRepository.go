package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"realty-avito/internal/models"
)

type HousesRepository interface {
	CreateHouse(ctx context.Context, houseModel models.CreateHouseEntity) (*models.House, error)
}

type housesRepository struct {
	pool *pgxpool.Pool
}

func NewHousesRepository(pool *pgxpool.Pool) HousesRepository {
	return &housesRepository{pool: pool}
}

func (r *housesRepository) CreateHouse(ctx context.Context, houseModel models.CreateHouseEntity) (*models.House, error) {
	query := squirrel.
		Insert("houses").
		PlaceholderFormat(squirrel.Dollar).
		Columns("address", "year", "developer").
		Values(houseModel.Address, houseModel.Year, houseModel.Developer).
		Suffix("RETURNING id, created_at, address, year, developer")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var house models.House

	err = r.pool.
		QueryRow(ctx, sql, args...).
		Scan(&house.ID, &house.CreatedAt, &house.Address, &house.Year, &house.Developer)
	if err != nil {
		return nil, err
	}

	return &house, nil
}
