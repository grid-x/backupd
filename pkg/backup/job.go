package backup

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

type DataStore interface {
	ExportTo(tmpdir string) (string, error)
}

type Storage interface {
	Copy(localfile string, remotefile string) error
}

type Config struct {
	// Name of backup this job belongs to
	name string

	tempDir       string
	tempDirPrefix string
}

func DefaultConfig() Config {
	return Config{
		tempDir:       "", // empty string forces to use system's defaults
		tempDirPrefix: "",
	}
}

type BackupJobStatus struct {
	Name     string
	Duration time.Duration
	Error    error
}

func errorStatus(name string, start time.Time, err error) BackupJobStatus {
	end := time.Now()
	return BackupJobStatus{
		Name:     name,
		Duration: end.Sub(start),
		Error:    err,
	}
}

type BackupJob struct {
	conf      Config
	datastore DataStore
	storage   Storage
	statusc   chan BackupJobStatus
}

func NewBackupJob(ds DataStore, s Storage, name string, statusc chan BackupJobStatus) *BackupJob {
	conf := DefaultConfig()
	conf.name = name

	return &BackupJob{
		datastore: ds,
		storage:   s,
		conf:      conf,
		statusc:   statusc,
	}
}

func (b *BackupJob) Run() {
	start := time.Now()

	tmpDir, err := ioutil.TempDir(b.conf.tempDir, b.conf.tempDirPrefix)
	if err != nil {
		b.statusc <- errorStatus(b.conf.name, start, err)
		return
	}

	localfile, err := b.datastore.ExportTo(tmpDir)
	if err != nil {
		b.statusc <- errorStatus(b.conf.name, start, err)
		return
	}

	if err := b.storage.Copy(localfile, filepath.Base(localfile)); err != nil {
		b.statusc <- errorStatus(b.conf.name, start, err)
		return
	}

	end := time.Now()

	b.statusc <- BackupJobStatus{
		Name:     b.conf.name,
		Duration: end.Sub(start),
	}
}
