package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      string        `yaml:"port"`
	Copyright string        `yaml:"copyright"`
	Admin     AdminConfig   `yaml:"admin"`
	Session   SessionConfig `yaml:"session"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SessionConfig struct {
	Secret string `yaml:"secret"`
}

var AppConfig Config

func LoadConfig() error {
	if err := loadConfigFile("config.yaml"); err == nil {
		overrideFromEnv()
		return nil
	}

	if err := loadConfigFile("config.demo.yaml"); err == nil {
		overrideFromEnv()
		return nil
	}

	AppConfig = Config{
		Port:      "8080",
		Copyright: "AI导航 © 2024",
		Admin: AdminConfig{
			Username: "admin",
			Password: "admin123",
		},
		Session: SessionConfig{
			Secret: "your-secret-key-here",
		},
	}
	overrideFromEnv()
	return nil
}

func loadConfigFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(&AppConfig)
}

func overrideFromEnv() {
	if port := os.Getenv("PORT"); port != "" {
		AppConfig.Port = port
	}

	if copyright := os.Getenv("COPYRIGHT"); copyright != "" {
		AppConfig.Copyright = copyright
	}

	if username := os.Getenv("ADMIN_USERNAME"); username != "" {
		AppConfig.Admin.Username = username
	}

	if password := os.Getenv("ADMIN_PASSWORD"); password != "" {
		AppConfig.Admin.Password = password
	}

	if secret := os.Getenv("SESSION_SECRET"); secret != "" {
		AppConfig.Session.Secret = secret
	}
}
