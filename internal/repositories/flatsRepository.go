package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"realty-avito/internal/models"
)

// TODO: разнести модели на слои: репо - entity, сервис - model и хэндлеры - DTO
type FlatsRepository interface {
	GetFlatsByHouseID(ctx context.Context, houseID int64) ([]models.Flat, error)
	GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]models.Flat, error)
}

type flatsRepository struct {
	pool *pgxpool.Pool
}

func NewFlatsRepository(pool *pgxpool.Pool) FlatsRepository {
	return &flatsRepository{pool: pool}
}

func (r *flatsRepository) GetFlatsByHouseID(ctx context.Context, houseID int64) ([]models.Flat, error) {
	query := squirrel.
		Select("id", "house_id", "price", "rooms", "status").
		From("flats").
		Where(squirrel.Eq{"house_id": houseID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flats []models.Flat

	for rows.Next() {
		var flat models.Flat

		err := rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			return nil, err
		}

		flats = append(flats, flat)
	}

	return flats, nil
}

func (r *flatsRepository) GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]models.Flat, error) {
	query := squirrel.
		Select("id", "house_id", "price", "rooms", "status").
		From("flats").
		Where(
			squirrel.Eq{"house_id": houseID,
				"status": models.StatusApproved}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flats []models.Flat
	for rows.Next() {
		var flat models.Flat
		if err := rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			return nil, err
		}
		flats = append(flats, flat)
	}

	return flats, nil
}
