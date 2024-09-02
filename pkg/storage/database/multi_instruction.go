package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

type MultiInstruction struct {
	db *MindPalaceDB
	tx *sql.Tx
	ID uuid.UUID
}

func NewMultiInstruction(db *MindPalaceDB) *MultiInstruction {
	return &MultiInstruction{
		db: db,
		ID: uuid.New(),
		tx: nil,
	}
}

func (t *MultiInstruction) CurrentSchema() string {
	return t.db.CurrentSchema
}

func (t *MultiInstruction) Begin() error {
	loggers.Log.TXInfo(context.Background(), t.ID, "BEGIN Transaction")
	tx, err := t.db.DB().Begin()
	if err != nil {
		return err
	}
	t.tx = tx
	return nil
}

func (t *MultiInstruction) Commit() error {
	loggers.Log.TXInfo(context.Background(), t.ID, "COMMIT Transaction")
	return t.tx.Commit()
}

func (t *MultiInstruction) Rollback() error {
	loggers.Log.TXInfo(context.Background(), t.ID, "ROLLBACK Transaction")
	return t.tx.Rollback()
}

func (t *MultiInstruction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(context.Background(), query, args...)
}

func (t *MultiInstruction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(context.Background(), query, args...)
}

func (t *MultiInstruction) Exec(query string, args ...interface{}) error {
	_, err := t.tx.ExecContext(context.Background(), query, args...)

	return err
}

func (t *MultiInstruction) Prepare(query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(context.Background(), query)
}
