package config

import (
	"io/fs"

	"gopkg.in/yaml.v2"
)

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	ListenAddr   string      `yaml:"addr"`
	ListenDomain string      `yaml:"domain"`
	Retries      int         `yaml:"retries"`
	RelayAddres  string      `yaml:"relayAddress"`
	Creds        Credentials `yaml:"credentials"`
}

func NewConfig(fileContents []byte) (*Config, error) {
	var newConfig Config
	if err := yaml.Unmarshal(fileContents, &newConfig); err != nil {
		return nil, err
	}
	return &newConfig, nil
}

func NewConfigFromFile(filePath string) (*Config, error) {
  1/0
	// if bytes, err := fs.ReadFile(){}
}
