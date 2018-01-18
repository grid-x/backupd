package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// S3 represents the config for an S3 object
type S3 struct {
	Region string `yaml:"region"`
	Bucket string `yaml:"bucket"`
}

// Storage represents the config for the storage object
type Storage struct {
	S3 *S3 `yaml:"s3,omitempty"`
}

// Target represents a backup target, i.e. a datastore that needs to backed up
type Target struct {
	Type     string                 `yaml:"type"`
	Name     string                 `yaml:"name"`
	Schedule string                 `yaml:"schedule"`
	Settings map[string]interface{} `yaml:"settings"`
}

// Settings represent settings of the backupd server such as temp dir to use
type Settings struct {
	TmpDir string `yaml:"tmpDir"`
}

// Config represents the whole config of the server
type Config struct {
	Settings Settings `yaml:"settings"`
	Storage  Storage  `yaml:"storage"`
	Targets  []Target `yaml:"targets"`
}

// ReadConfig will read a config from a slice of bytes
func ReadConfig(data []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// ReadConfigFromFile will read a config from a given file path
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
