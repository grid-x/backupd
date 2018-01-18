package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

// DataStore represents a datastore that can export data to a file
type DataStore interface {
	fmt.Stringer
	ExportTo(tmpdir string) (string, error)
}

// Storage represents a storage for a backup file and provides a mechanism to
// copy a localfile to the remote location
type Storage interface {
	Copy(localfile string, remotefile string) error
}

// Config is the configuration for a backup job
type Config struct {
	// Name of backup this job belongs to
	Name string

	TempDir       string
	TempDirPrefix string
}

// DefaultConfig returns the default config
func DefaultConfig() Config {
	return Config{
		TempDir:       "", // empty string forces to use system's defaults
		TempDirPrefix: "",
	}
}

// JobStatus represents the status, e.g. error and duration, of a
// finished backup job
type JobStatus struct {
	Name     string
	Duration time.Duration
	Error    error
}

func errorStatus(name string, start time.Time, err error) JobStatus {
	end := time.Now()
	return JobStatus{
		Name:     name,
		Duration: end.Sub(start),
		Error:    err,
	}
}

// Job represents a backup job and contains all settings for the job
type Job struct {
	conf      Config
	datastore DataStore
	storage   Storage
	statusc   chan JobStatus

	logger log.FieldLogger
}

// NewJob creates a new backup job for the given data store and the
// storage using the config provided. When executed the status will be written
// to the status channel
func NewJob(ds DataStore, s Storage, conf Config, statusc chan JobStatus) *Job {
	return &Job{
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

// Run will execute the backup job and is blocking
func (b *Job) Run() {
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

	b.statusc <- JobStatus{
		Name:     b.conf.Name,
		Duration: end.Sub(start),
	}

}
