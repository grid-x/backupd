package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type S3 struct {
	Region string `yaml:"region"`
	Bucket string `yaml:"bucket"`
}

type Storage struct {
	S3 *S3 `yaml:"s3,omitempty"`
}

type Target struct {
	Type     string                 `yaml:"type"`
	Name     string                 `yaml:"name"`
	Schedule string                 `yaml:"schedule"`
	Settings map[string]interface{} `yaml:"settings"`
}

type Settings struct {
	TmpDir string `yaml:"tmpDir"`
}

type Config struct {
	Settings Settings `yaml:"settings"`
	Storage  Storage  `yaml:"storage"`
	Targets  []Target `yaml:"targets"`
}

func ReadConfig(data []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func ReadConfigFromFile(filepath string) (*Config, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ReadConfig(b)
}
