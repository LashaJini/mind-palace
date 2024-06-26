package database

import (
	"database/sql"
	"fmt"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	_ "github.com/lib/pq"
)

type MindPalaceDB struct {
	db               *sql.DB
	ConnectionString string
	CurrentSchema    string
	cfg              *common.Config
}

func InitDB(cfg *common.Config) *MindPalaceDB {
	connStr := cfg.DBAddr()
	db, err := sql.Open(cfg.DB_DRIVER, connStr)
	errors.On(err).Exit()

	err = db.Ping()
	errors.On(err).Exit()

	return &MindPalaceDB{
		db:               db,
		ConnectionString: connStr,
		cfg:              cfg,
		CurrentSchema:    cfg.DB_DEFAULT_NAMESPACE,
	}
}

func (m *MindPalaceDB) DB() *sql.DB {
	return m.db
}

func (m *MindPalaceDB) ListMPSchemas() ([]string, error) {
	var results []string
	q := fmt.Sprintf("SELECT schema_name FROM information_schema.schemata WHERE schema_name LIKE '%%%s' OR schema_name = 'public'", m.cfg.DB_SCHEMA_SUFFIX)
	rows, err := m.db.Query(q)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name string
		)
		err := rows.Scan(&name)
		if err != nil {
			return results, err
		}

		results = append(results, name)
	}

	common.Log.Info().Msgf("found schemas %v", results)
	return results, nil
}

func (m *MindPalaceDB) CreateSchema(user string) error {
	schema := m.ConstructSchema(user)
	_, err := m.db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema))
	if err != nil {
		return err
	}

	common.Log.Info().Msgf("created schema '%s'", schema)
	m.SetSchema(schema)
	return nil
}

func (m *MindPalaceDB) SetSchema(schema string) {
	m.CurrentSchema = schema
}

func (m *MindPalaceDB) ConstructSchema(user string) string {
	return user + m.cfg.DB_SCHEMA_SUFFIX
}
