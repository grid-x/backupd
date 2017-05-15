package backup

import (
	"io/ioutil"
	"path/filepath"
)

type DataStore interface {
	ExportTo(tmpdir string) error
}

type Storage interface {
	Copy(localfile string, remotefile string) error
}

type Config struct {
	tempDir       string
	tempDirPrefix string
}

func DefaultConfig() Config {
	return Config{
		tempDir:       "", // empty string forces to use system's defaults
		tempDirPrefix: "",
	}
}

type Backup struct {
	conf      Config
	datastore DataStore
	storage   Storage
}

func NewBackup(config *Config, ds DataStore, s Storage) *Backup {
	if config == nil {
		tmp := DefaultConfig()
		config = &tmp
	}
	return &Backup{
		datastore: ds,
		storage:   s,
		conf:      *config,
	}
}

func (b *Backup) Run() error {
	tmpDir, err := ioutil.TempDir(b.conf.tempDir, b.conf.tempDirPrefix)
	if err != nil {
		return err
	}

	if err := b.datastore.ExportTo(tmpDir); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			// Expect files and not directories
			continue
		}

		basename := f.Name()
		localfile := filepath.Join(tmpDir, basename)
		if err := b.storage.Copy(localfile, basename); err != nil {
			return err
		}
	}
	return nil
}
