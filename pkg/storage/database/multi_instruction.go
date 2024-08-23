package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

type MultiInstruction struct {
	db  *MindPalaceDB
	tx  *sql.Tx
	ID  uuid.UUID
	ctx context.Context
}

func NewMultiInstruction(ctx context.Context, db *MindPalaceDB) *MultiInstruction {
	return &MultiInstruction{
		db:  db,
		ctx: ctx,
		ID:  uuid.New(),
		tx:  nil,
	}
}

func (t *MultiInstruction) CurrentSchema() string {
	return t.db.CurrentSchema
}

func (t *MultiInstruction) Begin() error {
	loggers.Log.TXInfo(t.ctx, t.ID, "BEGIN Transaction")
	tx, err := t.db.DB().BeginTx(t.ctx, nil)
	if err != nil {
		return err
	}
	t.tx = tx
	return nil
}

func (t *MultiInstruction) Commit() error {
	loggers.Log.TXInfo(t.ctx, t.ID, "COMMIT Transaction")
	return t.tx.Commit()
}

func (t *MultiInstruction) Rollback() error {
	loggers.Log.TXInfo(t.ctx, t.ID, "ROLLBACK Transaction")
	return t.tx.Rollback()
}

func (t *MultiInstruction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(t.ctx, query, args...)
}

func (t *MultiInstruction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(t.ctx, query, args...)
}

func (t *MultiInstruction) Exec(query string, args ...interface{}) error {
	_, err := t.tx.ExecContext(t.ctx, query, args...)

	return err
}

func (t *MultiInstruction) Prepare(query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(t.ctx, query)
}
