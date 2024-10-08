package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type Addon struct {
	Name        string
	Description string
	InputTypes  types.IOTypes
	OutputTypes types.IOTypes
	Response    *pb.AddonResponse
}

func (a *Addon) GetName() string {
	return a.Name
}

func (a *Addon) GetDescription() string {
	return a.Description
}

func (a *Addon) GetInputTypes() types.IOTypes {
	return a.InputTypes
}

func (a *Addon) GetOutputTypes() types.IOTypes {
	return a.OutputTypes
}

func (a *Addon) GetResponse() *pb.AddonResponse {
	return a.Response
}

func (a *Addon) SetResponse(response *pb.AddonResponse) {
	a.Response = response
}

func (a *Addon) Empty() bool {
	return a == nil || a.Name == ""
}

func (a Addon) Action(ctx context.Context, db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) error {
	return nil
}

func (a Addon) String() string {
	return fmt.Sprintf(`
Name: %s
Description: %s
Input types: %s
Output types: %s
Last response: %s
`, a.Name, a.Description, a.InputTypes, a.OutputTypes, a.Response)
}

func ToAddons(addonResult *pb.AddonResult) ([]types.IAddon, error) {
	var addons []types.IAddon
	if addonResult != nil {
		for key, addonResponse := range addonResult.Map {
			addon, err := Find(key)
			if err != nil {
				return nil, mperrors.On(err).Wrap("could not find addon")
			}

			if !addon.Empty() && addonResponse.Success {
				addon.SetResponse(addonResponse)

				addons = append(addons, addon)
			}
		}
	}

	return addons, nil
}

var List = []types.IAddon{
	&DefaultAddonInstance,
	&SummaryAddonInstance,
	&KeywordsAddonInstance,
}

func Find(name string) (types.IAddon, error) {
	for _, a := range List {
		if a.GetName() == name {
			return a, nil
		}
	}

	return nil, mperrors.Onf("addon %s not found", name)
}
