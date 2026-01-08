package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Admin   AdminConfig   `json:"admin"`
	Session SessionConfig `json:"session"`
}

type AdminConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionConfig struct {
	Secret string `json:"secret"`
}

var AppConfig Config

func LoadConfig() error {
	file, err := os.Open("config/config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&AppConfig)
}
