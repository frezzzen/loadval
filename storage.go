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

func (s *StorageAPI) SaveTemplates(templatesJSON string) error {
	return s.SaveData("templates.json", json.RawMessage(templatesJSON))
}

func (s *StorageAPI) LoadTemplates() (string, error) {
	return s.LoadData("templates.json")
}

func (s *StorageAPI) SaveAgentLoadouts(loadoutsJSON string) error {
	return s.SaveData("agent-loadouts.json", json.RawMessage(loadoutsJSON))
}

func (s *StorageAPI) LoadAgentLoadouts() (string, error) {
	return s.LoadData("agent-loadouts.json")
}

func (s *StorageAPI) SaveSettings(settingsJSON string) error {
	return s.SaveData("settings.json", json.RawMessage(settingsJSON))
}

func (s *StorageAPI) LoadSettings() (string, error) {
	return s.LoadData("settings.json")
}
