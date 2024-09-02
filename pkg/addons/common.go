package addons

import (
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

func revert(tx *database.MultiInstruction) {
	if r := recover(); r != nil {
		rollback(tx)
	}
}

func rollback(tx *database.MultiInstruction) {
	err := tx.Rollback()
	errors.On(err).PanicWithMsg("failed to rollback")
}
