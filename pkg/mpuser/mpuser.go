package mpuser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/config"
)

type Config struct {
	Config UserConfigRoot `json:"config"`
}

func (u *Config) Steps() []string {
	return u.Config.Text.Steps
}

func (u *Config) EnableAddon(addon addons.IAddon) error {
	inputTypes := addon.GetInputTypes()
	var needsUpdate bool

	for _, inputType := range inputTypes {
		if inputType == addons.Text {
			steps := u.Config.Text.Steps

			canAppend := true
			for _, step := range steps {
				if step == addon.GetName() {
					canAppend = false
					fmt.Printf("Addon '%s' is already enabled for '%s'.\n", addon.GetName(), addons.Text)
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

func (u *Config) DisableAddon(addon addons.IAddon) error {
	inputTypes := addon.GetInputTypes()
	var needsUpdate bool

	for _, inputType := range inputTypes {
		if inputType == addons.Text {
			steps := u.Config.Text.Steps

			for index, step := range steps {
				if step == addon.GetName() {
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

	return os.WriteFile(config.UserConfigPath(u.Config.User, true), d, 0777)
}

func NewUserConfig(user string) *Config {

	return &Config{
		Config: UserConfigRoot{
			User: user,
			Text: Input{
				Steps: []string{
					addons.Default,
				},
			},
		},
	}
}

func ReadUserConfig(user string) (*Config, error) {
	d, err := os.ReadFile(config.UserConfigPath(user, true))
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
