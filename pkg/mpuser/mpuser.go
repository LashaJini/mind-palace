package mpuser

import (
	"encoding/json"
	"os"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/types"
)

type Config struct {
	Config UserConfigRoot `json:"config"`
}

func (u *Config) Steps() []string {
	return u.Config.Text.Steps
}

func (u *Config) EnableAddon(addon types.IAddon) error {
	inputTypes := addon.GetInputTypes()
	var needsUpdate bool

	for _, inputType := range inputTypes {
		if inputType == types.Text {
			steps := u.Config.Text.Steps

			canAppend := true
			for _, step := range steps {
				if step == addon.GetName() {
					canAppend = false
					common.Log.Info().Msgf("addon '%s' is already enabled for '%s'.\n", addon.GetName(), types.Text)
					break
				}
			}

			if canAppend {
				u.Config.Text.Steps = append(steps, addon.GetName())
				needsUpdate = true
			}
		}
	}

	if needsUpdate {
		return u.Update()
	}

	return nil
}

func (u *Config) DisableAddon(addon types.IAddon) error {
	inputTypes := addon.GetInputTypes()
	var needsUpdate bool

	for _, inputType := range inputTypes {
		if inputType == types.Text {
			steps := u.Config.Text.Steps

			for index, step := range steps {
				addonName := addon.GetName()
				if step == addonName && addonName != types.AddonDefault {
					u.Config.Text.Steps = append(steps[:index], steps[index+1:]...)
					needsUpdate = true
					break
				}
			}
		}
	}

	if needsUpdate {
		return u.Update()
	}

	return nil
}

func (u *Config) Update() error {
	d, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return os.WriteFile(common.UserConfigPath(u.Config.User, true), d, 0777)
}

func CreateNewUser(user string) (*Config, error) {
	return CreateMindPalace(user)
}

func NewUserConfig(user string) *Config {
	return &Config{
		Config: UserConfigRoot{
			User: user,
			Text: Input{
				Steps: []string{
					types.AddonDefault,
				},
			},
		},
	}
}

func DeleteUser(user string) error {
	userpath, exists, err := UserExists(user)
	if err != nil {
		return err
	}

	if exists {
		err := os.RemoveAll(userpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func UserExists(user string) (string, bool, error) {
	mindPalaceUserPath := common.UserPath(user, true)
	exists, err := common.DirExists(mindPalaceUserPath)
	return mindPalaceUserPath, exists, err
}

func ReadConfig(user string) (*Config, error) {
	d, err := os.ReadFile(common.UserConfigPath(user, true))
	if err != nil {
		return &Config{}, err
	}

	var userCfg Config
	err = json.Unmarshal(d, &userCfg)

	return &userCfg, err
}

type UserConfigRoot struct {
	User string `json:"user"`
	Text Input  `json:"text"` // I don't like this
}

type Input struct {
	Steps []string `json:"steps"`
}
