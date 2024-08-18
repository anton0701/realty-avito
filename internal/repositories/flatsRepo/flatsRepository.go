package flatsRepo

import (
	"context"
	"realty-avito/internal/errors"

	"github.com/Masterminds/squirrel"

	"realty-avito/internal/client/db"
)

const (
	tableName = "flats"

	idColumn          = "id"
	houseIDColumn     = "house_id"
	priceColumn       = "price"
	roomsColumn       = "rooms"
	statusColumn      = "status"
	moderatorIDColumn = "moderator_id"
	createdAtColumn   = "created_at"
	updatedAtColumn   = "updated_at"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=FlatsRepository
type FlatsRepository interface {
	GetFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error)
	GetFlatByFlatID(ctx context.Context, flatID int64) (*FlatEntity, error)
	GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error)
	CreateFlat(ctx context.Context, flatModel CreateFlatEntity) (*FlatEntity, error)
	UpdateFlat(ctx context.Context, updateFlatModel UpdateFlatEntity) (*FlatEntity, error)
}

type flatsRepository struct {
	db db.Client
}

func NewFlatsRepository(db db.Client) FlatsRepository {
	return &flatsRepository{db: db}
}

func (r *flatsRepository) GetFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error) {
	selectBuilder := squirrel.
		Select(idColumn, houseIDColumn, priceColumn, roomsColumn, statusColumn).
		From(tableName).
		Where(squirrel.Eq{houseIDColumn: houseID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "flatsRepository.GetFlatsByHouseID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
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
	selectBuilder := squirrel.
		Select(idColumn, houseIDColumn, priceColumn, roomsColumn, statusColumn, moderatorIDColumn).
		From(tableName).
		Where(squirrel.Eq{idColumn: flatID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "flatsRepository.GetFlatByFlatID",
		QueryRaw: query,
	}

	var flat FlatEntity

	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status, &flat.ModeratorID)
	if err != nil {
		return nil, err
	}

	return &flat, nil
}

func (r *flatsRepository) GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]FlatEntity, error) {
	selectBuilder := squirrel.
		Select(idColumn, houseIDColumn, priceColumn, roomsColumn, statusColumn).
		From(tableName).
		Where(
			squirrel.Eq{
				houseIDColumn: houseID,
				statusColumn:  StatusApproved,
			}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "flatsRepository.GetApprovedFlatsByHouseID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
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
	insertBuilder := squirrel.
		Insert(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(houseIDColumn, priceColumn, roomsColumn, statusColumn).
		Values(flatEntity.HouseID, flatEntity.Price, flatEntity.Rooms, StatusCreated).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "flatsRepository.CreateFlat",
		QueryRaw: query,
	}

	var flatID int64

	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&flatID)
	if err != nil {
		if repo_errors.IsForeignKeyViolation(err) {
			return nil, &repo_errors.ErrHouseNotFound{HouseID: flatEntity.HouseID}
		}

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
	updateBuilder := squirrel.
		Update(tableName).
		Set(statusColumn, updateFlatEntity.Status).
		Set(moderatorIDColumn, updateFlatEntity.ModeratorID).
		Set(updatedAtColumn, updateFlatEntity.UpdatedAt).
		Where(squirrel.Eq{idColumn: updateFlatEntity.ID}).
		Suffix("RETURNING id, house_id, price, rooms, status, moderator_id").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "flatsRepository.UpdateFlat",
		QueryRaw: query,
	}

	var flat FlatEntity
	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status, &flat.ModeratorID)
	if err != nil {
		return nil, err
	}

	return &flat, nil
}
