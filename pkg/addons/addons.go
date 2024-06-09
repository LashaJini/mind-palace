package addons

import (
	"fmt"
	"strings"

	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
)

type Addon struct {
	Name        string
	Description string
	InputTypes  Types
	OutputTypes Types
	Output      any
	Action      func() (bool, error)
}

func (a *Addon) Empty() bool {
	return a == nil || a.Name == ""
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

func ToAddons(addonResult *pb.AddonResult) ([]Addon, error) {
	var addons []Addon
	if addonResult != nil {
		for key, value := range addonResult.Data {
			addon := Find(key)
			if !addon.Empty() && value.Success {
				addon.Output = strings.Join(value.Value, ", ")

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

var List = []Addon{
	{
		Name:        Default,
		Description: "Default",
		InputTypes:  []Type{Text},
		OutputTypes: []Type{Text},
		Action:      func() (bool, error) { return true, nil },
	},
	{
		Name:        ResourceSummary,
		Description: "Summarizes a resource",
		InputTypes:  []Type{Text},
		OutputTypes: []Type{Text},
		Action:      func() (bool, error) { return true, nil },
	},
	{
		Name:        ResourceKeywords,
		Description: "Extracts keywords from a resource",
		InputTypes:  []Type{Text},
		OutputTypes: []Type{Text},
		Action:      func() (bool, error) { return true, nil },
	},
}

func Find(name string) Addon {
	for _, a := range List {
		if a.Name == name {
			return a
		}
	}

	return Addon{}
}
