package flat

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FlatsRepository interface {
	GetFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error)
	GetFlatByFlatID(ctx context.Context, flatID int64) (*FlatEntity, error)
	GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error)
	CreateFlat(ctx context.Context, flatModel CreateFlatEntity) (*FlatEntity, error)
	UpdateFlat(ctx context.Context, updateFlatModel UpdateFlatEntity) (*FlatEntity, error)
}

type flatsRepository struct {
	pool *pgxpool.Pool
}

func NewFlatsRepository(pool *pgxpool.Pool) FlatsRepository {
	return &flatsRepository{pool: pool}
}

func (r *flatsRepository) GetFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error) {
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

	var flats []FlatEntity

	for rows.Next() {
		var flat FlatEntity

		err := rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			return nil, err
		}

		flats = append(flats, flat)
	}

	return flats, nil
}

func (r *flatsRepository) GetFlatByFlatID(ctx context.Context, flatID int64) (*FlatEntity, error) {
	query := squirrel.
		Select("id", "house_id", "price", "rooms", "status", "moderator_id").
		From("flats").
		Where(squirrel.Eq{"id": flatID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var flat FlatEntity

	err = r.pool.
		QueryRow(ctx, sql, args...).
		Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status, &flat.ModeratorID)
	if err != nil {
		return nil, err
	}

	return &flat, nil
}

func (r *flatsRepository) GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error) {
	query := squirrel.
		Select("id", "house_id", "price", "rooms", "status").
		From("flats").
		Where(
			squirrel.Eq{
				"house_id": houseID,
				"status":   StatusApproved,
			}).
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

	var flats []FlatEntity
	for rows.Next() {
		var flat FlatEntity
		if err := rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			return nil, err
		}
		flats = append(flats, flat)
	}

	return flats, nil
}

func (r *flatsRepository) CreateFlat(ctx context.Context, flatEntity CreateFlatEntity) (*FlatEntity, error) {
	query := squirrel.
		Insert("flats").
		PlaceholderFormat(squirrel.Dollar).
		Columns("house_id", "price", "rooms", "status").
		Values(flatEntity.HouseID, flatEntity.Price, flatEntity.Rooms, StatusCreated).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var flatID int64

	err = r.pool.
		QueryRow(ctx, sql, args...).
		Scan(&flatID)
	if err != nil {
		return nil, err
	}

	var flat = &FlatEntity{
		ID:      flatID,
		HouseID: flatEntity.HouseID,
		Price:   flatEntity.Price,
		Rooms:   flatEntity.Rooms,
		Status:  StatusCreated,
	}

	return flat, nil
}

func (r *flatsRepository) UpdateFlat(ctx context.Context, updateFlatEntity UpdateFlatEntity) (*FlatEntity, error) {
	query := squirrel.
		Update("flats").
		Set("status", updateFlatEntity.Status).
		Set("moderator_id", updateFlatEntity.ModeratorID).
		Where(squirrel.Eq{"id": updateFlatEntity.ID}).
		Suffix("RETURNING id, house_id, price, rooms, status, moderator_id").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var flat FlatEntity
	err = r.pool.
		QueryRow(ctx, sql, args...).
		Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status, &flat.ModeratorID)
	if err != nil {
		return nil, err
	}

	return &flat, nil
}
