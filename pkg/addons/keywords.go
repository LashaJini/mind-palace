package addons

import (
	"fmt"

	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type KeywordsAddon struct {
	Addon
}

func (k *KeywordsAddon) Action(db *database.MindPalaceDB) (bool, error) {
	keywords := k.Output
	fmt.Println(keywords)
	return true, nil
}

func (k *KeywordsAddon) SetOutput(output any) {
	k.Output = output
}

var KeywordsAddonInstance = KeywordsAddon{
	Addon: Addon{
		Name:        ResourceKeywords,
		Description: "Extracts keywords from a resource",
		InputTypes:  []Type{Text},
		OutputTypes: []Type{Text},
	},
}
