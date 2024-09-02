package types

import (
	"context"
	"strings"

	"github.com/google/uuid"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type IAddon interface {
	Empty() bool
	Action(ctx context.Context, db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) error
	String() string

	GetName() string
	GetDescription() string
	GetInputTypes() IOTypes
	GetOutputTypes() IOTypes
	GetResponse() *pb.AddonResponse

	SetResponse(response *pb.AddonResponse)
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
