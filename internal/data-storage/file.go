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

const PROFILE_FILE_NAME string = ".cstc_profiles"
const STATE_FILE_NAME string = "state"

type systemState struct {
	CurrentProfile string `json:"active_profile"`
}

type UserProfile struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var pathSeparator, profileStatePath string
var state *systemState
var dataState *profileData

func init() {

	if os.Getenv("GOOS") == "windows" {
		pathSeparator = "\\"
		// TODO: Fix this path for windows
		profileStatePath = "./etc/testsource"
	} else {
		pathSeparator = "/"
		// TODO: Fix this to use root path
		profileStatePath = "./etc/testsource"
	}

	profileFile := fmt.Sprintf("%s%s%s", ".", pathSeparator, PROFILE_FILE_NAME)
	stateFile := fmt.Sprintf("%s%s%s", ".", pathSeparator, STATE_FILE_NAME)

	for _, file := range []string{stateFile, profileFile} {
		if _, err := os.Stat(file); errors.Is(err, fs.ErrNotExist) {
			f, err := os.Create(file)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var emptyDoc map[string]UserProfile = make(map[string]UserProfile)
			err = json.NewEncoder(f).Encode(&emptyDoc)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if _, err := os.Stat(profileStatePath); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(profileStatePath, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	state = new(systemState)
	loadState()
	dataState = NewProfileData(state.CurrentProfile)
	loadActiveProfileDataState()
}

// File ops
func getStoredProfiles() (map[string]UserProfile, error) {
	profileFile := fmt.Sprintf("%s%s%s", ".", pathSeparator, PROFILE_FILE_NAME)

	data, err := os.ReadFile(profileFile)
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

func saveProfiles(userProfiles map[string]UserProfile) error {
	stateFile := fmt.Sprintf("%s%s%s", ".", pathSeparator, PROFILE_FILE_NAME)

	data, err := json.Marshal(userProfiles)
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0644)
}

func loadState() error {
	data, err := os.ReadFile(fmt.Sprintf("%s%s%s", ".", pathSeparator, STATE_FILE_NAME))

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, state)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func loadActiveProfileDataState() error {
	dataFile, err := os.ReadFile(fmt.Sprintf("%s%s%s", profileStatePath, pathSeparator, GetActiveProfileName()))

	if err != nil {
		return err
	}

	err = json.Unmarshal(dataFile, &dataState)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func saveState() error {
	stateFile := fmt.Sprintf("%s%s%s", ".", pathSeparator, STATE_FILE_NAME)

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0644)
}

func createProfileStateFile(profileName string) error {
	stateFile := fmt.Sprintf("%s%s%s", profileStatePath, pathSeparator, profileName)
	if _, err := os.Stat(stateFile); err != nil {
		file, err := os.Create(stateFile)
		if err != nil {
			return err
		}
		defer file.Close()

		var emptymap map[string]string = make(map[string]string)
		err = json.NewEncoder(file).Encode(emptymap)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("file '%s' already exists", stateFile)
}

// Profile File Operations
func GetActiveUserProfile() (*UserProfile, error) {
	userProfiles, err := getStoredProfiles()
	if err != nil {
		return nil, err
	}
	profileName := GetActiveProfileName()

	for key, userProfile := range userProfiles {
		if key == profileName {
			return &userProfile, nil
		}
	}
	return nil, errors.New("user profile not found")
}

func CreateNewUserProfile(profileName string, profile UserProfile) error {
	userProfiles, err := getStoredProfiles()
	if err != nil {
		return err
	}

	userProfiles[profileName] = profile
	if state.CurrentProfile == "" {
		state.CurrentProfile = profileName
		if err = saveState(); err != nil {
			return err
		}
	}

	if err = createProfileStateFile(profileName); err != nil {
		return err
	}
	return saveProfiles(userProfiles)
}

func DeleteUserProfile(profileName string) error {
	userProfiles, err := getStoredProfiles()
	if err != nil {
		return err
	}

	delete(userProfiles, profileName)

	if len(userProfiles) == 0 {
		state.CurrentProfile = ""
		if err = clearActiveProfileState(); err != nil {
			return err
		}
	}

	if err = os.Remove(fmt.Sprintf("%s%s%s", profileStatePath, pathSeparator, profileName)); err != nil {
		return err
	}

	return saveProfiles(userProfiles)
}

func GetUserProfileList() ([]string, error) {
	userProfiles, err := getStoredProfiles()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(userProfiles))
	for key := range userProfiles {
		keys = append(keys, key)
	}

	return keys, nil
}

func IsEmailUsedInProfile(email string) (bool, error) {
	userProfiles, err := getStoredProfiles()
	if err != nil {
		return false, err
	}

	for _, userProfile := range userProfiles {
		if userProfile.Email == email {
			return true, nil
		}
	}
	return false, nil
}

// State File Operations
func SetNewActiveProfileState(profile string) error {
	userProfiles, err := getStoredProfiles()
	if err != nil {
		return err
	}

	for k := range userProfiles {
		if k == profile {
			state.CurrentProfile = k
			saveState()
			return nil
		}
	}
	return fmt.Errorf("user profile '%s' not found", profile)
}

func clearActiveProfileState() error {
	state.CurrentProfile = ""
	saveState()
	return nil
}

func GetActiveProfileName() string {
	return state.CurrentProfile
}

func GetAvailableCoursesForActiveProfile() ([]Course, error) {
	return dataState.Courses, nil
}

func GetAvailableAssignmentsForCourse(courseId int) ([]Assignment, error) {
	for _, course := range dataState.Courses {
		if course.Id == courseId {
			return course.Assignments, nil
		}
	}
	return nil, fmt.Errorf("course-id '%d' not found", courseId)
}

func GetAllAssignmentsForActiveProfile() ([]Assignment, error) {
	var assignments []Assignment
	for _, course := range dataState.Courses {
		assignments = append(assignments, course.Assignments...)
	}
	return assignments, nil
}

func GetAssignmentCourse(assignmentName string) (Course, error) {
	return dataState.GetAssignmentCourse(assignmentName)
}

// Profile Data State File Operations
