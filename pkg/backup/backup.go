package backup

import (
	"io/ioutil"
	"path/filepath"
)

type DataStore interface {
	ExportTo(tmpdir string) (string, error)
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

type BackupJob struct {
	conf      Config
	datastore DataStore
	storage   Storage
}

func NewBackupJob(config *Config, ds DataStore, s Storage) *BackupJob {
	if config == nil {
		tmp := DefaultConfig()
		config = &tmp
	}
	return &BackupJob{
		datastore: ds,
		storage:   s,
		conf:      *config,
	}
}

func (b *BackupJob) Run() error {
	tmpDir, err := ioutil.TempDir(b.conf.tempDir, b.conf.tempDirPrefix)
	if err != nil {
		return err
	}

	localfile, err := b.datastore.ExportTo(tmpDir)
	if err != nil {
		return err
	}

	if err := b.storage.Copy(localfile, filepath.Base(localfile)); err != nil {
		return err
	}
	return nil
}
