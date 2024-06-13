package types

import (
	"strings"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type IAddon interface {
	Empty() bool
	Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) error
	String() string

	GetName() string
	GetDescription() string
	GetInputTypes() IOTypes
	GetOutputTypes() IOTypes
	GetOutput() any

	SetOutput(output any)
}

type IOType string
type IOTypes []IOType

const (
	Text IOType = "text"
)

func (t IOTypes) String() string {
	var result []string

	for _, ts := range t {
		result = append(result, string(ts))
	}

	return strings.Join(result, ", ")
}

var (
	AddonDefault          = "mind-palace-default"
	AddonResourceSummary  = "mind-palace-resource-summary"
	AddonResourceKeywords = "mind-palace-resource-keywords"
)
