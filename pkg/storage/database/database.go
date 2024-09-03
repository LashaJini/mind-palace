package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	_ "github.com/lib/pq"
)

const RETRY_COUNT = 3

type MindPalaceDB struct {
	db               *sql.DB
	ConnectionString string
	CurrentSchema    string
	cfg              *common.Config
}

func InitDB(cfg *common.Config) *MindPalaceDB {
	ctx := context.Background()
	connStr := cfg.DBAddr()
	loggers.Log.Debug(ctx, connStr)

	var err error
	db, err := sql.Open(cfg.DB_DRIVER, connStr)
	mperrors.On(err).Panic()

	m := &MindPalaceDB{
		db:               db,
		ConnectionString: connStr,
		cfg:              cfg,
		CurrentSchema:    cfg.DB_DEFAULT_NAMESPACE,
	}

	if err := m.Ping(ctx); err != nil {
		loggers.Log.Fatal(ctx, err, "")
		panic(err)
	}

	return m
}

func (m *MindPalaceDB) Ping(ctx context.Context) error {
	var err error
	for i := 1; i <= RETRY_COUNT; i++ {
		err = m.db.Ping()
		if err != nil {
			loggers.Log.Warn(ctx, "database ping '%d' failed (retrying in 1 sec), reason: %v", i, err)
		} else {
			loggers.Log.Info(ctx, "database ping '%d' successful", i)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

func (m *MindPalaceDB) DB() *sql.DB {
	return m.db
}

func (m *MindPalaceDB) ListMPSchemas() ([]string, error) {
	ctx := context.Background()
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

	loggers.Log.Info(ctx, "found schemas %v", results)
	return results, nil
}

func (m *MindPalaceDB) CreateSchema(user string) (string, error) {
	ctx := context.Background()
	schema := m.ConstructSchema(user)
	_, err := m.db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema))
	if err != nil {
		return "", err
	}

	loggers.Log.Info(ctx, "created schema '%s'", schema)
	m.SetSchema(schema)
	return schema, nil
}

func (m *MindPalaceDB) SetSchema(schema string) {
	m.CurrentSchema = schema
}

func (m *MindPalaceDB) ConstructSchema(user string) string {
	return user + m.cfg.DB_SCHEMA_SUFFIX
}
