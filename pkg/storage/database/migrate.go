package database

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

type InMemorySource struct {
	Migrations *source.Migrations
}

func (s *InMemorySource) Open(url string) (source.Driver, error) {
	return s, nil
}

func (s *InMemorySource) Close() error {
	return nil
}

func (s *InMemorySource) First() (uint, error) {
	v, ok := s.Migrations.First()
	if !ok {
		return 0, &os.PathError{Op: "first", Path: "", Err: os.ErrNotExist}
	}
	return v, nil
}

func (s *InMemorySource) Prev(version uint) (uint, error) {
	v, ok := s.Migrations.Prev(version)
	if !ok {
		return 0, &os.PathError{Op: "prev", Path: "", Err: os.ErrNotExist}
	}
	return v, nil
}

func (s *InMemorySource) Next(version uint) (uint, error) {
	v, ok := s.Migrations.Next(version)
	if !ok {
		return 0, &os.PathError{Op: "next", Path: "", Err: os.ErrNotExist}
	}
	return v, nil
}

func (s *InMemorySource) ReadUp(version uint) (io.ReadCloser, string, error) {
	if m, ok := s.Migrations.Up(version); ok {
		return io.NopCloser(bytes.NewReader([]byte(m.Raw))), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: "readUp", Path: "", Err: os.ErrNotExist}
}

func (s *InMemorySource) ReadDown(version uint) (io.ReadCloser, string, error) {
	if m, ok := s.Migrations.Down(version); ok {
		return io.NopCloser(bytes.NewReader([]byte(m.Raw))), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: "readDown", Path: "", Err: os.ErrNotExist}
}

func NewInMemorySource() *InMemorySource {
	return &InMemorySource{
		Migrations: source.NewMigrations(),
	}
}

// TODO: when ups (steps) exceeds total number of migrations (or total steps - last migration version),
// calculate the difference and apply
func Up(cfg *common.Config, sqlTemplates []common.SQLTemplate) (*InMemorySource, int) {
	ctx := context.Background()
	inMemorySource := NewInMemorySource()
	ups := 0

	ext := ".up.sql"
	err := filepath.Walk(cfg.MIGRATIONS_DIR, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			migrationID := migrationIDFromFile(path)

			var sqlBuffer bytes.Buffer
			for _, sqlTemplate := range sqlTemplates {
				err := sqlTemplate.Inject(&sqlBuffer, path)
				mperrors.On(err).Exit()
			}

			migration := &source.Migration{
				Version:    migrationID,
				Direction:  source.Up,
				Raw:        sqlBuffer.String(),
				Identifier: path,
			}
			inMemorySource.Migrations.Append(migration)

			// count total files ending with *.up.sql in migrations dir
			// ls -l migrations | grep ".up.sql" | wc -l
			ups++
		}
		return nil
	})

	mperrors.On(err).Exit()

	loggers.Log.Info(ctx, "total 'up' migrations found: %d", ups)

	return inMemorySource, ups
}

func Down(cfg *common.Config, sqlTemplates []common.SQLTemplate) *InMemorySource {
	inMemorySource := NewInMemorySource()

	ext := ".down.sql"
	err := filepath.Walk(cfg.MIGRATIONS_DIR, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			migrationID := migrationIDFromFile(path)

			var sqlBuffer bytes.Buffer
			for _, sqlTemplate := range sqlTemplates {
				err := sqlTemplate.Inject(&sqlBuffer, path)
				mperrors.On(err).Exit()
			}

			migration := &source.Migration{
				Version:    migrationID,
				Direction:  source.Down,
				Raw:        sqlBuffer.String(),
				Identifier: path,
			}
			inMemorySource.Migrations.Append(migration)
		}
		return nil
	})

	mperrors.On(err).Exit()

	return inMemorySource
}

func migrationIDFromFile(path string) uint {
	fileName := filepath.Base(path)
	timestamp := strings.Split(fileName, "_")[0]
	migrationID, err := strconv.Atoi(timestamp)
	mperrors.On(err).Exit()

	return uint(migrationID)
}

func CommitMigration(cfg *common.Config, inMemorySource *InMemorySource, steps int) error {
	mm, err := migrate.NewWithSourceInstance("in-memory", inMemorySource, cfg.DBAddr())
	if err != nil {
		return mperrors.On(err).Wrap("could not create migration instance")
	}

	return migrationSteps(mm, steps)
}

func migrationSteps(m *migrate.Migrate, steps int) error {
	err := m.Steps(steps)
	if err == migrate.ErrNoChange {
		loggers.Log.Warn(context.Background(), "no migrations applied. %s", err)
		return nil
	}

	if err != nil {
		return mperrors.On(err).Wrap("migration steps failed")
	}

	return nil
}
