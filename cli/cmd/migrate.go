package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database structures.",
	Long:  "Migrate database structures. Create new tables, columns, indexes and so on.",
	Run:   Migrate,
}

var (
	UP      int
	DOWN    int
	VERSION bool
	FORCE   int
	CREATE  string
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().IntVarP(&UP, "up", "u", 1, "migrate up by <n>. Setting to 0 migrates all")
	migrateCmd.PersistentFlags().IntVarP(&DOWN, "down", "d", 1, "migrate down by <n>")
	migrateCmd.Flags().BoolVarP(&VERSION, "version", "v", false, "last migration version")
	migrateCmd.PersistentFlags().IntVarP(&FORCE, "force", "f", 1, "set migration version. Don't run migrations")
	migrateCmd.Flags().StringVarP(&CREATE, "create", "c", "", "create migration")

	migrateCmd.MarkFlagsOneRequired("up", "down", "version", "force", "create")
	migrateCmd.MarkFlagsMutuallyExclusive("up", "down", "version", "force", "create")
}

func Migrate(cmd *cobra.Command, args []string) {
	cfg := common.NewConfig()
	m, err := migrate.New("file://"+cfg.MIGRATIONS_DIR, cfg.DBAddr())
	errors.On(err).Exit()

	if VERSION {
		version, dirty, err := m.Version()
		errors.On(err).Exit()
		common.Log.Info().Msgf("version: %d dirty: %t", version, dirty)

		return
	}

	if cmd.Flags().Changed("force") {
		err := m.Force(FORCE)
		errors.On(err).Exit()
		return
	}

	if cmd.Flags().Changed("create") {
		timestamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s", timestamp, CREATE)
		upFilePath := filepath.Join(cfg.MIGRATIONS_DIR, fmt.Sprintf("%s.up.sql", fileName))
		downFilePath := filepath.Join(cfg.MIGRATIONS_DIR, fmt.Sprintf("%s.down.sql", fileName))

		upFile, err := os.Create(upFilePath)
		errors.On(err).Exit()
		defer upFile.Close()

		downFile, err := os.Create(downFilePath)
		errors.On(err).Exit()
		defer downFile.Close()

		return
	}

	// up or down migrations

	db := database.InitDB(cfg)
	defer db.DB().Close()

	schemas, err := db.ListMPSchemas()
	errors.On(err).Exit()

	sqlTemplateDatas := []SQLTemplateData{}
	for _, schema := range schemas {
		sqlTemplateData := SQLTemplateData{
			Namespace: schema,
		}

		sqlTemplateDatas = append(sqlTemplateDatas, sqlTemplateData)
	}

	steps := 0
	var inMemorySource *InMemorySource
	if cmd.Flags().Changed("down") {
		inMemorySource = down(cfg, sqlTemplateDatas)

		steps = -DOWN
	} else {
		var ups int
		inMemorySource, ups = up(cfg, sqlTemplateDatas)

		steps = UP
		if steps == 0 {
			steps = ups
		}
	}

	mm, err := migrate.NewWithSourceInstance("in-memory", inMemorySource, cfg.DBAddr())
	errors.On(err).Exit()

	migrationSteps(mm, steps)
}

func migrationSteps(m *migrate.Migrate, steps int) {
	err := m.Steps(steps)
	if err == migrate.ErrNoChange {
		common.Log.Warn().Msgf("no migrations applied. %s", err)
		return
	}
	errors.On(err).Exit()

	common.Log.Info().Msg("successfully migrated")
}

type SQLTemplateData struct {
	Namespace string
}

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
func up(cfg *common.Config, sqlTemplateDatas []SQLTemplateData) (*InMemorySource, int) {
	inMemorySource := NewInMemorySource()
	ups := 0

	ext := ".up.sql"
	err := filepath.Walk(cfg.MIGRATIONS_DIR, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			migrationID := migrationIDFromFile(path)

			var sqlBuffer bytes.Buffer
			for _, sqlTemplateData := range sqlTemplateDatas {
				inject(&sqlBuffer, path, sqlTemplateData)
			}
			migration := &source.Migration{
				Version:    uint(migrationID),
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

	errors.On(err).Exit()

	common.Log.Info().Msgf("total 'up' migrations found: %d", ups)

	return inMemorySource, ups
}

func down(cfg *common.Config, sqlTemplateDatas []SQLTemplateData) *InMemorySource {
	inMemorySource := NewInMemorySource()

	ext := ".down.sql"
	err := filepath.Walk(cfg.MIGRATIONS_DIR, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			migrationID := migrationIDFromFile(path)

			var sqlBuffer bytes.Buffer
			for _, sqlTemplateData := range sqlTemplateDatas {
				inject(&sqlBuffer, path, sqlTemplateData)
			}
			migration := &source.Migration{
				Version:    uint(migrationID),
				Direction:  source.Down,
				Raw:        sqlBuffer.String(),
				Identifier: path,
			}
			inMemorySource.Migrations.Append(migration)
		}
		return nil
	})

	errors.On(err).Exit()

	return inMemorySource
}

func migrationIDFromFile(path string) uint {
	fileName := filepath.Base(path)
	timestamp := strings.Split(fileName, "_")[0]
	migrationID, err := strconv.Atoi(timestamp)
	errors.On(err).Exit()

	return uint(migrationID)
}

func inject(sqlBuffer *bytes.Buffer, path string, sqlTemplateData SQLTemplateData) {
	f, err := os.ReadFile(path)
	errors.On(err).Exit()

	tmpl, err := template.New("sql").Parse(string(f))
	errors.On(err).Exit()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, sqlTemplateData)
	errors.On(err).Exit()

	sqlBuffer.WriteString(buf.String())
}
