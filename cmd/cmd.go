package main

import "gopkg.in/yaml.v3"

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	RelayAddres string      `yaml:"realayAddress"`
	Creds       Credentials `yaml:"credentials"`
}

func NewConfig(fileContents []byte) (*Config, error) {
	var newConfig Config
	if err := yaml.Unmarshal(fileContents, &newConfig); err != nil {
		return nil, err
	}
	return &newConfig, nil
}
