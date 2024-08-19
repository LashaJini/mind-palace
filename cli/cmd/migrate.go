package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
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

	sqlTemplates := common.NewSQLTemplates(schemas)

	steps := 0
	var inMemorySource *database.InMemorySource
	if cmd.Flags().Changed("down") {
		inMemorySource = database.Down(cfg, sqlTemplates)

		steps = -DOWN
	} else {
		var ups int
		inMemorySource, ups = database.Up(cfg, sqlTemplates)

		steps = UP
		if steps == 0 {
			steps = ups
		}
	}

	database.CommitMigration(cfg, inMemorySource, steps)
}
