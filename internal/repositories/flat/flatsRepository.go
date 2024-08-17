package flat

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"

	"realty-avito/internal/client/db"
)

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
		Select("id", "house_id", "price", "rooms", "status").
		From("flats").
		Where(squirrel.Eq{"house_id": houseID}).
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
		Select("id", "house_id", "price", "rooms", "status", "moderator_id").
		From("flats").
		Where(squirrel.Eq{"id": flatID}).
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
		Select("id", "house_id", "price", "rooms", "status").
		From("flats").
		Where(
			squirrel.Eq{
				"house_id": houseID,
				"status":   StatusApproved,
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
	log.Printf("flatsRepository CreateFlat context: %v", ctx)

	insertBuilder := squirrel.
		Insert("flats").
		PlaceholderFormat(squirrel.Dollar).
		Columns("house_id", "price", "rooms", "status").
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
		Update("flats").
		Set("status", updateFlatEntity.Status).
		Set("moderator_id", updateFlatEntity.ModeratorID).
		Set("updated_at", updateFlatEntity.UpdatedAt).
		Where(squirrel.Eq{"id": updateFlatEntity.ID}).
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
