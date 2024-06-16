package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
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

	steps := 0
	// TODO: annoying error when there are no more migrations to apply
	if cmd.Flags().Changed("up") {
		ups := 0
		if UP == 0 {
			// count total files ending with *.up.sql in migrations dir
			// ls -l migrations | grep ".up.sql" | wc -l
			ext := ".up.sql"
			err := filepath.Walk(cfg.MIGRATIONS_DIR, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && strings.HasSuffix(path, ext) {
					ups++
				}
				return nil
			})
			errors.On(err).Exit()

			common.Log.Info().Msgf("total 'up' migrations: %d", ups)
		}
		steps = ups
	} else if cmd.Flags().Changed("down") {
		steps = -DOWN
	}

	migrationSteps(m, steps)
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
