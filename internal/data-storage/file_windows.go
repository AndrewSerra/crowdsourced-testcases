//go:build windows
// +build windows

/*
 * Created on Tue Dec 17 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package datastorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

const PATH_SEPARATOR string = "\\"
const CONFIG_FILE_NAME string = ".cstc"

type UserProfile struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

func init() {
	// configFile := fmt.Sprintf("%s%s%s", "$HOME", PATH_SEPARATOR, CONFIG_FILE_NAME)
	configFile := fmt.Sprintf("%s%s%s", ".", PATH_SEPARATOR, CONFIG_FILE_NAME)

	if _, err := os.Stat(configFile); errors.Is(err, fs.ErrNotExist) {
		fmt.Println("config file does not exist")
		os.Create(configFile)
	}
}

func getConfig() (map[string]UserProfile, error) {
	// configFile := fmt.Sprintf("%s%s%s", "$HOME", PATH_SEPARATOR, CONFIG_FILE_NAME)
	configFile := fmt.Sprintf("%s%s%s", ".", PATH_SEPARATOR, CONFIG_FILE_NAME)

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var userProfiles map[string]UserProfile
	err = json.Unmarshal(data, &userProfiles)
	if err != nil {
		return nil, err
	}

	return userProfiles, nil
}

func saveConfig(userProfiles map[string]UserProfile) error {
	// configFile := fmt.Sprintf("%s%s%s", "$HOME", PATH_SEPARATOR, CONFIG_FILE_NAME)
	configFile := fmt.Sprintf("%s%s%s", ".", PATH_SEPARATOR, CONFIG_FILE_NAME)

	data, err := json.Marshal(userProfiles)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

func GetUserProfile(profileName string) (*UserProfile, error) {
	userProfiles, err := getConfig()
	if err != nil {
		return nil, err
	}

	for key, userProfile := range userProfiles {
		if key == profileName {
			return &userProfile, nil
		}
	}
	return nil, errors.New("user profile not found")
}

func CreateNewUserProfile(profileName string, profile UserProfile) error {
	userProfiles, err := getConfig()
	if err != nil {
		return err
	}

	userProfiles[profileName] = profile
	return saveConfig(userProfiles)
}

func DeleteUserProfile(profileName string) error {
	userProfiles, err := getConfig()
	if err != nil {
		return err
	}

	delete(userProfiles, profileName)
	return saveConfig(userProfiles)
}
