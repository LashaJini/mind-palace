package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lashajini/mind-palace/config"
	_ "github.com/lib/pq"
)

type MindPalaceDB struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) *MindPalaceDB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.DB_USER, cfg.DB_PASS, cfg.DB_NAME, cfg.DB_PORT)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &MindPalaceDB{db}
}
