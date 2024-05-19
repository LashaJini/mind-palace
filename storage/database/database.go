package database

import (
	"database/sql"
	"log"

	"github.com/lashajini/mind-palace/config"
	_ "github.com/lib/pq"
)

type MindPalaceDB struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) *MindPalaceDB {
	db, err := sql.Open("postgres", cfg.DBAddr())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &MindPalaceDB{db}
}
