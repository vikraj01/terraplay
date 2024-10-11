package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveConfig[T any](data T, baseDir, filename string) error {
	configPath := filepath.Join(baseDir, filename)

	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write to config file: %v", err)
	}

	return nil
}

func LoadConfig[T any](baseDir, filename string, data *T) error {
	configPath := filepath.Join(baseDir, filename)

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}
