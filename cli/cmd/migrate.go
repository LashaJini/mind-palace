package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database structures.",
	Long:  "Migrate database structures. Create new tables, columns, indexes and so on.",
	Run:   Migrate,
}

var (
	MIGRATE_UP      int
	MIGRATE_DOWN    int
	MIGRATE_VERSION bool
	MIGRATE_FORCE   int
	MIGRATE_CREATE  string
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().IntVarP(&MIGRATE_UP, "up", "u", 1, "migrate up by <n>. Setting to 0 migrates all")
	migrateCmd.PersistentFlags().IntVarP(&MIGRATE_DOWN, "down", "d", 1, "migrate down by <n>")
	migrateCmd.Flags().BoolVarP(&MIGRATE_VERSION, "version", "v", false, "last migration version")
	migrateCmd.PersistentFlags().IntVarP(&MIGRATE_FORCE, "force", "f", 1, "set migration version. Don't run migrations")
	migrateCmd.Flags().StringVarP(&MIGRATE_CREATE, "create", "c", "", "create migration")

	migrateCmd.MarkFlagsOneRequired("up", "down", "version", "force", "create")
	migrateCmd.MarkFlagsMutuallyExclusive("up", "down", "version", "force", "create")
}

func Migrate(cmd *cobra.Command, args []string) {
	cfg := common.NewConfig()
	m, err := migrate.New("file://"+cfg.MIGRATIONS_DIR, cfg.DBAddr())
	mperrors.On(err).Exit()

	if MIGRATE_VERSION {
		version, dirty, err := m.Version()
		mperrors.On(err).ExitWithMsgf("version: %d dirty: %t", version, dirty)

		fmt.Println(version)
		return
	}

	if cmd.Flags().Changed("force") {
		err := m.Force(MIGRATE_FORCE)
		mperrors.On(err).Exit()
		return
	}

	if cmd.Flags().Changed("create") {
		timestamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s", timestamp, MIGRATE_CREATE)
		upFilePath := filepath.Join(cfg.MIGRATIONS_DIR, fmt.Sprintf("%s.up.sql", fileName))
		downFilePath := filepath.Join(cfg.MIGRATIONS_DIR, fmt.Sprintf("%s.down.sql", fileName))

		upFile, err := os.Create(upFilePath)
		mperrors.On(err).Exit()
		defer _close(upFile)

		downFile, err := os.Create(downFilePath)
		mperrors.On(err).Exit()
		defer _close(downFile)

		return
	}

	// up or down migrations

	db := database.InitDB(cfg)
	defer closeDB(db)

	schemas, err := db.ListMPSchemas()
	mperrors.On(err).Exit()

	sqlTemplates := common.NewSQLTemplates(schemas)

	steps := 0
	var inMemorySource *database.InMemorySource
	if cmd.Flags().Changed("down") {
		inMemorySource = database.Down(cfg, sqlTemplates)

		steps = -MIGRATE_DOWN
	} else {
		var ups int
		inMemorySource, ups = database.Up(cfg, sqlTemplates)

		steps = MIGRATE_UP
		if steps == 0 {
			steps = ups
		}
	}

	err = database.CommitMigration(cfg, inMemorySource, steps)
	mperrors.On(err).Exit()
}

func _close(f *os.File) {
	if err := f.Close(); err != nil {
		mperrors.On(err).Exit()
	}
}
