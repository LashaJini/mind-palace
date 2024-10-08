package addons

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type SummaryAddon struct {
	Addon
}

func (s *SummaryAddon) Action(ctx context.Context, db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	summary := s.Response.GetSummaryResponse().Summary

	tx := database.NewMultiInstruction(db)
	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				err = mperrors.On(rollbackErr).Wrap("summary rollback failed")
			} else {
				err = mperrors.Onf("(recovered) panic: %v", r)
			}
		}
	}()

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				err = mperrors.On(rollbackErr).Wrap("summary rollback failed")
			}
		}
	}()

	err = tx.Begin()
	if err != nil {
		return mperrors.On(err).Wrap("summary transaction begin failed")
	}

	summaryID := uuid.New()
	select {
	case memoryID := <-memoryIDC:
		err = models.InsertSummaryTx(ctx, tx, memoryID, summaryID, summary)
		if err != nil {
			return mperrors.On(err).Wrap("failed to insert summary")
		}

		err = tx.Commit(ctx)
		if err != nil {
			return mperrors.On(err).Wrap("summary transaction commit failed")
		}

		return nil
	case <-ctx.Done():
		if ctx.Err() != nil {
			return mperrors.On(ctx.Err()).Wrap("summary addon action failed")
		}

		return nil
	}
}

var SummaryAddonInstance = SummaryAddon{
	Addon: Addon{
		Name:        types.AddonResourceSummary,
		Description: "Summarizes a resource",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
