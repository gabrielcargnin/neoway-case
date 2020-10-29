package db

import (
	"context"
	"neoway-case/errors"
	"neoway-case/schema"
)

// Repository is an interface allowing the user to encapsulate the data access logic.
// This way we ensure that our business layer don't need to know how the data is being persisted
type Repository interface {
	Close()
	InsertConsumption(ctx context.Context, row []schema.Consumption) error
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertConsumption(ctx context.Context, rows []schema.Consumption) error {
	err := impl.InsertConsumption(ctx, rows)
	if err != nil {
		return errors.E(err, errors.Op("db.InsertConsumption"))
	}
	return err
}
