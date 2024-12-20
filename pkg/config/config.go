package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

const dbPathKey = "DBPath"
const defaultDbPath = "pogreb.db"

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Connection struct {
	ListenAddr        string `yaml:"addr"`
	ListenDomain      string `yaml:"domain"`
	AllowInsecureAuth bool   `yaml:"insecure_auth"`
}

type Config struct {
	Connection    Connection  `yaml:"connection"`
	Retries       int         `yaml:"retries"`
	RelayAddres   string      `yaml:"relayAddress"`
	Creds         Credentials `yaml:"credentials"`
	AllowInsecure bool        `yaml:"allow_insecure"`
}

type DatabaseConfig struct {
	FilePath string `yaml:"path"`
}

func NewConfig(fileContents []byte) (*Config, error) {
	var newConfig Config
	if err := yaml.Unmarshal(fileContents, &newConfig); err != nil {
		return nil, err
	}
	return &newConfig, nil
}

func NewConfigFromFile(filePath string) (*Config, error) {
	if bytes, err := os.ReadFile(filePath); err != nil {
		return nil, err
	} else {
		return NewConfig(bytes)
	}
}

func NewDbConfigFromEnv() *DatabaseConfig {
	if val, ok := os.LookupEnv(dbPathKey); !ok {
		return &DatabaseConfig{FilePath: defaultDbPath}
	} else {
		return &DatabaseConfig{FilePath: val}
	}
}
