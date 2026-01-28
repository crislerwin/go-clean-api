package postgres

import (
	"errors"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/lib/pq"
)

func TranslateError(err error) error {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23503":
			return mapForeignKeyError(pqErr)
		case "23505":
			return entity.ErrEventSoldOut
		}

	}

	return err
}

func mapForeignKeyError(err *pq.Error) error {
	switch err.Constraint {
	case "orders_user_id_fkey":
		return usecase.ErrUserNotFound
	case "orders_event_id_fkey":
		return usecase.ErrEventNotFound
	default:
		return errors.New("conflict: foreign key violation")
	}

}
