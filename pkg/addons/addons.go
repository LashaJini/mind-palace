package addons

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type IAddon interface {
	Empty() bool
	Action(db *database.MindPalaceDB, memoryID uuid.UUID, args ...any) (bool, error)
	String() string

	GetName() string
	GetDescription() string
	GetInputTypes() Types
	GetOutputTypes() Types
	GetOutput() any

	SetOutput(output any)
}

type Addon struct {
	Name        string
	Description string
	InputTypes  Types
	OutputTypes Types
	Output      any
}

func (a *Addon) GetName() string {
	return a.Name
}

func (a *Addon) GetDescription() string {
	return a.Description
}

func (a *Addon) GetInputTypes() Types {
	return a.InputTypes
}

func (a *Addon) GetOutputTypes() Types {
	return a.OutputTypes
}

func (a *Addon) GetOutput() any {
	return a.Output
}

func (a *Addon) SetOutput(output any) {
	a.Output = output
}

func (a *Addon) Empty() bool {
	return a == nil || a.Name == ""
}

func (a Addon) Action(db *database.MindPalaceDB, memoryID uuid.UUID, args ...any) (bool, error) {
	return true, nil
}

func (a Addon) String() string {
	return fmt.Sprintf(`
Name: %s
Description: %s
Input types: %s
Output types: %s
Last output: %s
`, a.Name, a.Description, a.InputTypes, a.OutputTypes, a.Output)
}

func ToAddons(addonResult *pb.AddonResult) ([]IAddon, error) {
	var addons []IAddon
	if addonResult != nil {
		for key, value := range addonResult.Data {
			addon := Find(key)
			if !addon.Empty() && value.Success {
				addon.SetOutput(value.Value)

				addons = append(addons, addon)
			}
		}
	}

	return addons, nil
}

type Type string
type Types []Type

const (
	Text Type = "text"
)

func (t Types) String() string {
	var result []string

	for _, ts := range t {
		result = append(result, string(ts))
	}

	return strings.Join(result, ", ")
}

var (
	Default          = "mind-palace-default"
	ResourceSummary  = "mind-palace-resource-summary"
	ResourceKeywords = "mind-palace-resource-keywords"
)

var DefaultAddonInstance = Addon{
	Name:        Default,
	Description: "Default",
	InputTypes:  []Type{Text},
	OutputTypes: []Type{Text},
}

var SummaryAddonInstance = Addon{
	Name:        ResourceSummary,
	Description: "Summarizes a resource",
	InputTypes:  []Type{Text},
	OutputTypes: []Type{Text},
}

var List = []IAddon{
	&DefaultAddonInstance,
	&SummaryAddonInstance,
	&KeywordsAddonInstance,
}

func Find(name string) IAddon {
	for _, a := range List {
		if a.GetName() == name {
			return a
		}
	}

	return &Addon{}
}
