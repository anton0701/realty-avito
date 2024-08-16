package house

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HousesRepository interface {
	CreateHouse(ctx context.Context, createHouseEntity CreateHouseEntity) (*HouseEntity, error)
}

type housesRepository struct {
	pool *pgxpool.Pool
}

func NewHousesRepository(pool *pgxpool.Pool) HousesRepository {
	return &housesRepository{pool: pool}
}

func (r *housesRepository) CreateHouse(ctx context.Context, createHouseEntity CreateHouseEntity) (*HouseEntity, error) {
	query := squirrel.
		Insert("houses").
		PlaceholderFormat(squirrel.Dollar).
		Columns("address", "year", "developer").
		Values(createHouseEntity.Address, createHouseEntity.Year, createHouseEntity.Developer).
		Suffix("RETURNING id, created_at, address, year, developer")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var house HouseEntity

	err = r.pool.
		QueryRow(ctx, sql, args...).
		Scan(&house.ID, &house.CreatedAt, &house.Address, &house.Year, &house.Developer)
	if err != nil {
		return nil, err
	}

	return &house, nil
}
