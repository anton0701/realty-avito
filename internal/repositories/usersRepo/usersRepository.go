package usersRepo

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"realty-avito/internal/client/db"
)

const (
	usersTable             = "users"
	userIDColumn           = "id"
	userEmailColumn        = "email"
	userPasswordHashColumn = "password_hash"
	userTypeColumn         = "user_type"
	createdAtColumn        = "created_at"
	userUUIDColumn         = "uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user UserEntity) (*UserEntity, error)
	GetUserByCredentials(ctx context.Context, cred UserCredentials) (*UserEntity, error)
}

type userRepository struct {
	db db.Client
}

func NewUserRepository(db db.Client) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user UserEntity) (*UserEntity, error) {
	uuid := uuid.New().String()
	insertBuilder := squirrel.
		Insert(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns(userEmailColumn, userPasswordHashColumn, userTypeColumn, userUUIDColumn).
		Values(user.Email, user.PasswordHash, user.UserType, uuid).
		Suffix("RETURNING " + userIDColumn + ", " + createdAtColumn + ", " + userUUIDColumn)

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "userRepository.CreateUser",
		QueryRaw: query,
	}

	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&user.ID, &user.CreatedAt, &user.UUID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Код ошибки уникального констрэйнта
				if strings.Contains(pgErr.ConstraintName, "email") {
					return nil, ErrEmailExists
				}
			}
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByCredentials(ctx context.Context, cred UserCredentials) (*UserEntity, error) {
	log.Printf("userRepository.GetUserByCredentials cred.ID: %s", cred.ID)
	log.Printf("userRepository.GetUserByCredentials cred.PasswordHash: %s", cred.PasswordHash)

	builder := squirrel.
		Select(userIDColumn, userEmailColumn, userPasswordHashColumn, userTypeColumn, createdAtColumn).
		From(usersTable).
		Where(squirrel.Eq{
			userUUIDColumn: cred.ID,
		}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "userRepository.GetUserByCredentials",
		QueryRaw: query,
	}

	var user UserEntity
	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.UserType, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
