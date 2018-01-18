package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type DataStore interface {
	fmt.Stringer
	ExportTo(tmpdir string) (string, error)
}

type Storage interface {
	Copy(localfile string, remotefile string) error
}

type Config struct {
	// Name of backup this job belongs to
	Name string

	TempDir       string
	TempDirPrefix string
}

func DefaultConfig() Config {
	return Config{
		TempDir:       "", // empty string forces to use system's defaults
		TempDirPrefix: "",
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

	logger log.FieldLogger
}

func NewBackupJob(ds DataStore, s Storage, conf Config, statusc chan BackupJobStatus) *BackupJob {
	return &BackupJob{
		datastore: ds,
		storage:   s,
		conf:      conf,
		statusc:   statusc,

		logger: log.New().WithFields(log.Fields{
			"component": "backup-job",
			"datastore": ds,
			"name":      conf.Name,
		}),
	}
}

func (b *BackupJob) Run() {
	start := time.Now()

	b.logger.Info("Backup job started")
	defer b.logger.Info("Backup job finished")

	tmpDir, err := ioutil.TempDir(b.conf.TempDir, b.conf.TempDirPrefix)
	if err != nil {
		b.statusc <- errorStatus(b.conf.Name, start, err)
		return
	}

	localfile, err := b.datastore.ExportTo(tmpDir)
	if err != nil {
		b.statusc <- errorStatus(b.conf.Name, start, err)
		return
	}

	remotefile := filepath.Join(
		b.conf.Name,
		time.Now().Format(time.RFC3339),
		filepath.Base(localfile))
	if err := b.storage.Copy(localfile, remotefile); err != nil {
		b.statusc <- errorStatus(b.conf.Name, start, err)
		return
	}
	err = os.Remove(localfile)
	if err != nil {
		b.statusc <- errorStatus(b.conf.Name, start, err)
		return
	}

	end := time.Now()

	b.statusc <- BackupJobStatus{
		Name:     b.conf.Name,
		Duration: end.Sub(start),
	}

}
