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
	// 优先加载 config.yaml
	if err := loadConfigFile("config.yaml"); err == nil {
		return nil
	}

	// 然后加载 config.demo.yaml
	if err := loadConfigFile("config.demo.yaml"); err == nil {
		return nil
	}

	// 如果都不存在，使用默认配置
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
