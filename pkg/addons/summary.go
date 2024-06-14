package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/models"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc/client"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type SummaryAddon struct {
	Addon
}

func (s *SummaryAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	summary := s.Output.([]string)[0]
	rpcClient := args[0].(*rpcclient.Client)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())

	defer func() {
		if err != nil {
			fmt.Println(err)

			if err := tx.Rollback(); err != nil {
				panic(fmt.Errorf("Failed to rollback: %w", err))
			}
		}
	}()

	err = tx.Begin()
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}

	summaryID := uuid.New()
	memoryID := <-memoryIDC
	err = models.InsertSummaryTx(tx, memoryID, summaryID, summary)
	if err != nil {
		return fmt.Errorf("Failed to insert summary: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	err = rpcClient.VDBInsert(ctx, memoryID, summary)
	if err != nil {
		return fmt.Errorf("Failed to insert in vdb: %w", err)
	}

	return nil
}

var SummaryAddonInstance = SummaryAddon{
	Addon: Addon{
		Name:        types.AddonResourceSummary,
		Description: "Summarizes a resource",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
