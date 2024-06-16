package cli

import (
	"os"
	"path/filepath"
	"strings"

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
	UP   int
	DOWN int
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().IntVarP(&UP, "up", "u", 1, "migrate up by <n>. Setting to 0 migrates all.")
	migrateCmd.PersistentFlags().IntVarP(&DOWN, "down", "d", 1, "migrate down by <n>")

	migrateCmd.MarkFlagsOneRequired("up", "down")
	migrateCmd.MarkFlagsMutuallyExclusive("up", "down")
}

func Migrate(cmd *cobra.Command, args []string) {
	cfg := common.NewConfig()
	m, err := migrate.New("file://"+cfg.MIGRATIONS_DIR, cfg.DBAddr())
	errors.On(err).Exit()

	steps := 0
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
