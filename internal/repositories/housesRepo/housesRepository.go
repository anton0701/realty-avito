package housesRepo

import (
	"context"

	"github.com/Masterminds/squirrel"

	"realty-avito/internal/client/db"
)

const (
	tableName = "houses"

	idColumn        = "id"
	addressColumn   = "address"
	yearColumn      = "year"
	developerColumn = "developer"
	updatedAtColumn = "updated_at"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=HousesRepository
type HousesRepository interface {
	CreateHouse(ctx context.Context, createHouseEntity CreateHouseEntity) (*HouseEntity, error)
	UpdateHouseUpdatedAt(ctx context.Context, houseID int64) error
}

type housesRepository struct {
	db db.Client
}

func NewHousesRepository(db db.Client) HousesRepository {
	return &housesRepository{db: db}
}

func (r *housesRepository) CreateHouse(ctx context.Context, createHouseEntity CreateHouseEntity) (*HouseEntity, error) {
	insertBuilder := squirrel.
		Insert(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(addressColumn, yearColumn, developerColumn).
		Values(createHouseEntity.Address, createHouseEntity.Year, createHouseEntity.Developer).
		Suffix("RETURNING id, created_at, address, year, developer")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "housesRepository.CreateHouse",
		QueryRaw: query,
	}

	var house HouseEntity

	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&house.ID, &house.CreatedAt, &house.Address, &house.Year, &house.Developer)
	if err != nil {
		return nil, err
	}

	return &house, nil
}

func (r *housesRepository) UpdateHouseUpdatedAt(ctx context.Context, houseID int64) error {
	updateBuilder := squirrel.
		Update(tableName).
		Set(updatedAtColumn, squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{idColumn: houseID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "housesRepository.UpdateHouseUpdatedAt",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
