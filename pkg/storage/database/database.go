package database

import (
	"database/sql"
	"log"

	"github.com/lashajini/mind-palace/pkg/config"
	_ "github.com/lib/pq"
)

type MindPalaceDB struct {
	db               *sql.DB
	ConnectionString string
}

func InitDB(cfg *config.Config) *MindPalaceDB {
	connStr := cfg.DBAddr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &MindPalaceDB{db: db, ConnectionString: connStr}
}

func (m *MindPalaceDB) DB() *sql.DB {
	return m.db
}
