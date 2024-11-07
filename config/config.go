package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type RestServerConfig struct {
	Port int `yaml:"port"`
}

type AppConfig struct {
	Database     DatabaseConfig   `yaml:"database"`
	ServerConfig RestServerConfig `yaml:"restServerConfig"`
}

func LoadConfiguration(pathOfYaml string) (*AppConfig, error) {
	file, err := os.Open(pathOfYaml)
	if err != nil {
		log.Fatalf("Failed to find yaml file")
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read yaml file")
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal yaml file")
		return nil, err
	}

	return &config, nil
}
