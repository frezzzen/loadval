package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type StorageAPI struct{}

func NewStorageAPI() *StorageAPI {
	return &StorageAPI{}
}

func (s *StorageAPI) getAppDataDir() (string, error) {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		return "", fmt.Errorf("LOCALAPPDATA environment variable is not set")
	}

	appDir := filepath.Join(appData, "LOADVAL")

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}

	return appDir, nil
}

func (s *StorageAPI) getUserDataDir(userID string) (string, error) {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		return "", fmt.Errorf("LOCALAPPDATA environment variable is not set")
	}

	if userID == "" {
		return "", fmt.Errorf("userID is required")
	}

	userDir := filepath.Join(appData, "LOADVAL", "users", userID)

	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	return userDir, nil
}

func (s *StorageAPI) SaveData(filename string, data interface{}) error {
	appDir, err := s.getAppDataDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(appDir, filename)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *StorageAPI) LoadData(filename string) (string, error) {
	appDir, err := s.getAppDataDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(appDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "{}", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}

func (s *StorageAPI) SaveUserData(userID string, filename string, data interface{}) error {
	userDir, err := s.getUserDataDir(userID)
	if err != nil {
		return err
	}

	filePath := filepath.Join(userDir, filename)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	LogInfo(fmt.Sprintf("Saved %s for user: %s", filename, userID[:8]+"..."))
	return nil
}

func (s *StorageAPI) LoadUserData(userID string, filename string) (string, error) {
	userDir, err := s.getUserDataDir(userID)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(userDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		LogInfo(fmt.Sprintf("No existing %s for user: %s", filename, userID[:8]+"..."))
		return "{}", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	LogInfo(fmt.Sprintf("Loaded %s for user: %s", filename, userID[:8]+"..."))
	return string(data), nil
}

func (s *StorageAPI) SaveTemplates(userID string, templatesJSON string) error {
	return s.SaveUserData(userID, "templates.json", json.RawMessage(templatesJSON))
}

func (s *StorageAPI) LoadTemplates(userID string) (string, error) {
	return s.LoadUserData(userID, "templates.json")
}

func (s *StorageAPI) SaveAgentLoadouts(userID string, loadoutsJSON string) error {
	return s.SaveUserData(userID, "agent-loadouts.json", json.RawMessage(loadoutsJSON))
}

func (s *StorageAPI) LoadAgentLoadouts(userID string) (string, error) {
	return s.LoadUserData(userID, "agent-loadouts.json")
}

func (s *StorageAPI) SaveSettings(userID string, settingsJSON string) error {
	return s.SaveUserData(userID, "settings.json", json.RawMessage(settingsJSON))
}

func (s *StorageAPI) LoadSettings(userID string) (string, error) {
	return s.LoadUserData(userID, "settings.json")
}
