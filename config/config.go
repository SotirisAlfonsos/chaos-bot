package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PeerToken string `yaml:"peer_token,omitempty"`
	CertFile  string `yaml:"cert_file,omitempty"`
	KeyFile   string `yaml:"key_file,omitempty"`
}

func GetConfig(file string) (*Config, error) {
	return unmarshalConfFromFile(file)
}

func unmarshalConfFromFile(file string) (*Config, error) {
	config := &Config{}

	if file != "" {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			err = errors.Wrap(err, "could not read yml")
			return nil, err
		}

		if err = yaml.Unmarshal(yamlFile, &config); err != nil {
			err = errors.Wrap(err, "could not unmarshal yml")
			return nil, err
		}
	}

	return config, nil
}
