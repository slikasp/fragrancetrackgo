package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const configFileName = "appconfig.json"

type Config struct {
	LocalDbURL  string    `json:"local_db_url"`
	RemoteDbURL string    `json:"remote_db_url"`
	CurrentUser uuid.UUID `json:"current_user"`
}

type configJSON struct {
	LocalDbURL  string `json:"local_db_url"`
	RemoteDbURL string `json:"remote_db_url"`
	CurrentUser string `json:"current_user"`
}

func getConfigFilePath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(currentDir, configFileName)
	return path, nil
}

func Read() (Config, error) {
	// Read the config file in user's HOME directory
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Parse and return the Config struct
	decoder := json.NewDecoder(file)
	raw := configJSON{}
	err = decoder.Decode(&raw)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		LocalDbURL:  raw.LocalDbURL,
		RemoteDbURL: raw.RemoteDbURL,
		CurrentUser: uuid.Nil,
	}
	if strings.TrimSpace(raw.CurrentUser) != "" {
		if parsed, parseErr := uuid.Parse(raw.CurrentUser); parseErr == nil {
			cfg.CurrentUser = parsed
		}
	}

	return cfg, nil
}
