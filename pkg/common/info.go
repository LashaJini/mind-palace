package common

import (
	"encoding/json"
	"fmt"
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

	return os.WriteFile(InfoPath(true), d, 0777)
}

func CurrentUser() (string, error) {
	d, err := os.ReadFile(InfoPath(true))
	if err != nil {
		return "", err
	}

	info := MindPalaceInfo{}
	err = json.Unmarshal(d, &info)
	if err != nil {
		return "", err
	}

	if info.CurrentUser == "" {
		return "", fmt.Errorf("current user is not set")
	}

	return info.CurrentUser, err
}
