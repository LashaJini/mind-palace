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

func (t *MultiInstruction) BeginTx(ctx context.Context) error {
	loggers.Log.TXInfo(ctx, t.ID, "BEGIN Transaction")
	tx, err := t.db.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	t.tx = tx
	return nil
}

func (t *MultiInstruction) Commit(ctx context.Context) error {
	loggers.Log.TXInfo(ctx, t.ID, "COMMIT Transaction")
	return t.tx.Commit()
}

func (t *MultiInstruction) Rollback(ctx context.Context) error {
	loggers.Log.TXInfo(ctx, t.ID, "ROLLBACK Transaction")
	return t.tx.Rollback()
}

func (t *MultiInstruction) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *MultiInstruction) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *MultiInstruction) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.tx.ExecContext(ctx, query, args...)

	return err
}

func (t *MultiInstruction) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, query)
}
