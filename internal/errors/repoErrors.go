package repo_errors

import (
	"fmt"

	"github.com/jackc/pgconn"
)

type ErrHouseNotFound struct {
	HouseID int64
}

func (e *ErrHouseNotFound) Error() string {
	return fmt.Sprintf("house with ID %d not found", e.HouseID)
}

// Функция проверкяет что ошибка = нарушение внешнего ключа
func IsForeignKeyViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23503" // SQLSTATE 23503 - нарушение внешнего ключа
	}
	return false
}
