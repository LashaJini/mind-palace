package addons

import (
	"fmt"
	"strings"
)

type Addon struct {
	Name        string
	Description string
	Input       Types
	Output      Types
}

func (a Addon) String() string {
	return fmt.Sprintf(`Name: %s
Description: %s
Input: %s
Output: %s`, a.Name, a.Description, a.Input, a.Output)
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
	ResourceSummary = "mind-palace-resource-summary"
)

var List = []Addon{
	{Name: ResourceSummary, Description: "Summarizes a resource", Input: []Type{Text}, Output: []Type{Text}},
}

func Find(name string) Addon {
	for _, a := range List {
		if a.Name == name {
			return a
		}
	}

	return Addon{}
}
