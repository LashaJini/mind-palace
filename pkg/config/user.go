package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/addons"
)

type Type addons.Type

type UserConfig struct {
	Config UserConfigRoot `json:"config"`
}

func (u *UserConfig) Steps() []string {
	return u.Config.Text.Steps
}

func (u *UserConfig) EnableAddon(addon addons.Addon) error {
	inputTypes := addon.Input
	var needsUpdate bool

	for _, inputType := range inputTypes {
		if inputType == addons.Text {
			steps := u.Config.Text.Steps

			canAppend := true
			for _, step := range steps {
				if step == addon.Name {
					canAppend = false
					fmt.Printf("Addon '%s' is already enabled for '%s'.\n", addon.Name, addons.Text)
					break
				}
			}

			if canAppend {
				u.Config.Text.Steps = append(steps, addon.Name)
				needsUpdate = true
			}
		}
	}

	if needsUpdate {
		return u.Update()
	}

	return nil
}

func (u *UserConfig) DisableAddon(addon addons.Addon) error {
	inputTypes := addon.Input
	var needsUpdate bool

	for _, inputType := range inputTypes {
		if inputType == addons.Text {
			steps := u.Config.Text.Steps

			for index, step := range steps {
				if step == addon.Name {
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

func (u *UserConfig) Update() error {
	d, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return os.WriteFile(MindPalaceUserConfigPath(u.Config.User, true), d, 0777)
}

func NewUserConfig(user string) *UserConfig {
	return &UserConfig{
		Config: UserConfigRoot{
			User: user,
		},
	}
}

func ReadUserConfig(user string) (*UserConfig, error) {
	d, err := os.ReadFile(MindPalaceUserConfigPath(user, true))
	if err != nil {
		return &UserConfig{}, err
	}

	var userCfg UserConfig
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
