package database

import (
	"database/sql"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	_ "github.com/lib/pq"
)

type MindPalaceDB struct {
	db               *sql.DB
	ConnectionString string
}

func InitDB(cfg *common.Config) *MindPalaceDB {
	connStr := cfg.DBAddr()
	db, err := sql.Open(cfg.DB_DRIVER, connStr)
	errors.On(err).Exit()

	err = db.Ping()
	errors.On(err).Exit()

	return &MindPalaceDB{db: db, ConnectionString: connStr}
}

func (m *MindPalaceDB) DB() *sql.DB {
	return m.db
}
