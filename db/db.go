package db

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"neoway-case/errors"
	"neoway-case/schema"
)

// Postgres implementation of Repository interface.
// 2 implemented functions.

type PostgresRepository struct {
	db *sql.DB
}

// Return a new database connection
func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

// Close closes the database and prevents new queries from starting. But still waits for all already started queries to finish.
func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		log.Fatal(err)
	}
}

// InsertConsumption receives a consumption slice and persists all consumptions with COPY FROM statement for performance purposes
func (r *PostgresRepository) InsertConsumption(ctx context.Context, consumptions []schema.Consumption) error {
	const op errors.Op = "db.InsertConsumption"
	const kind errors.Kind = "Postgres error"
	txn, err := r.db.Begin()
	if err != nil {
		return errors.E(op, errors.Message("Could not start a postgres transaction"), kind)
	}

	stmt, err := txn.Prepare(pq.CopyIn("consumption", "cpf", "private", "incompleto", "data_ultima_compra", "ticket_medio", "ticket_ultima_compra", "loja_frequente", "loja_ultima_compra"))
	if err != nil {
		return errors.E(op, errors.Message("Could not create a postgres prepared statement"), kind)
	}

	for _, c := range consumptions {
		_, err = stmt.Exec(c.CPF, c.Private, c.Incompleto, c.DataUltimaCompra, c.TicketMedio, c.TicketUltimaCompra, c.LojaFrequente, c.LojaUltimaCompra)
		if err != nil {
			return errors.E(op, errors.Message("Could not execute postgres statement"), kind)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.E(op, errors.Message("Could not execute postgres statement"), kind)
	}

	err = stmt.Close()
	if err != nil {
		return errors.E(op, errors.Message("Could not close postgres statement"), kind)
	}

	err = txn.Commit()
	if err != nil {
		return errors.E(op, errors.Message("Could not commit transaction"), kind)
	}

	return err
}
