package config

import (
	"encoding/json"
	"os"
)

type MindPalaceInfo struct {
	CurrentUser string `json:"current_user"`
}

func UpdateMindPalaceInfo(info MindPalaceInfo) error {
	d, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return os.WriteFile(MindPalaceInfoPath(true), d, 0777)
}

func CurrentUser() (string, error) {
	d, err := os.ReadFile(MindPalaceInfoPath(true))
	if err != nil {
		return "", err
	}

	info := MindPalaceInfo{}
	err = json.Unmarshal(d, &info)

	return info.CurrentUser, err
}
