package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	influxd = "influxd"
)

func mkInfluxdCmd(database, host, tmpdir string, since *time.Time) *exec.Cmd {
	args := []string{"backup", "-host", host,
		"-database", database}

	if since != nil {
		args = append(args, "-since", (*since).Format(time.RFC3339))
	}

	args = append(args, tmpdir)

	return exec.Command(influxd, args...)
}

// Influx represents the datastore interface for the influx database
type Influx struct {
	endpoint string
	database string
	last     *time.Duration

	logger log.FieldLogger
}

// NewInflux creates a new influx object from the given settings
func NewInflux(endpoint, database string, last *time.Duration) *Influx {
	return &Influx{
		endpoint: endpoint,
		database: database,
		last:     last,

		logger: log.New().WithFields(log.Fields{
			"datastore": "influx",
		}),
	}
}

// String returns a string representation of the datastore
func (i *Influx) String() string {
	return "Influx"
}

// ExportTo exports the database contents to a file and uses the given tempdir
// required to satisfy the datastore interface
func (i *Influx) ExportTo(tmpdir string) (string, error) {
	var since *time.Time
	if i.last != nil {
		since = new(time.Time)
		*since = time.Now().Add(-1 * *i.last)
	}
	cmd := mkInfluxdCmd(i.database, i.endpoint, tmpdir, since)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error: %+v, Output: %s", string(out))
	}

	fs, err := ioutil.ReadDir(tmpdir)
	if err != nil {
		return "", err
	}

	backupFiles := []string{}
	for _, f := range fs {
		backupFiles = append(backupFiles, filepath.Join(tmpdir, f.Name()))
	}

	archive, err := ioutil.TempFile(tmpdir, "influxdb-archive-")
	if err != nil {
		return "", err
	}
	defer archive.Close()

	err = createTarArchive(archive, backupFiles)
	if err != nil {
		return "", err
	}

	for _, f := range backupFiles {
		if err := os.Remove(f); err != nil {
			i.logger.Warnf("Got error while removing %s: %+v", f, err)
		}
	}

	return archive.Name(), nil
}
