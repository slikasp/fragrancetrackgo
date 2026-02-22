package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

const configFileName = "appconfig.json"

type Config struct {
	UserDbURL      string    `json:"user_db_url"`
	FragranceDbURL string    `json:"fragrance_db_url"`
	CurrentUser    uuid.UUID `json:"current_user"`
}

type State struct {
	Users      *database.Queries
	Fragrances *database.Queries
	Cfg        *Config
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
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) SetUser(user uuid.UUID) error {

	// Set the user variable to input
	c.CurrentUser = user
	// Update the config file
	j, err := json.Marshal(c)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, j, 0777)
	if err != nil {
		return err
	}
	return nil
}
