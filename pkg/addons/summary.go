package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/models"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type SummaryAddon struct {
	Addon
}

func (s *SummaryAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	summary := s.Response.GetSummaryResponse().Summary

	vdbGrpcClient := args[0].(*vdbrpc.VDBGrpcClient)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db)
	defer revert(tx)

	err = tx.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	summaryID := uuid.New()
	memoryID := <-memoryIDC
	err = models.InsertSummaryTx(tx, memoryID, summaryID, summary)
	if err != nil {
		return fmt.Errorf("failed to insert summary: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err = vdbGrpcClient.Insert(ctx, []uuid.UUID{memoryID}, []string{summary})
	if err != nil {
		return fmt.Errorf("failed to insert in vdb: %w", err)
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
