package datastore

import (
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const (
	pgDump = "pg_dump"
)

// Postgres represents the datastore implementation for the postgres database
type Postgres struct {
	logger log.FieldLogger
	url    string
}

// NewPostgres creates a new postgres instance from the given connect URL
func NewPostgres(url string) *Postgres {
	return &Postgres{
		logger: log.New().WithFields(log.Fields{
			"component": "datastore",
			"datastore": "postgres",
		}),
		url: url,
	}
}

// String returns a string representation of the datastore
func (p *Postgres) String() string {
	return "postgresql"
}

// ExportTo exports the database contents to a file and uses the given tempdir
// required to satisfy the datastore interface
func (p *Postgres) ExportTo(tmpdir string) (string, error) {
	cmd := exec.Command(pgDump, p.url)
	f, err := ioutil.TempFile(tmpdir, "postgres-")
	if err != nil {
		return "", err
	}
	cmd.Stdout = f
	if err := cmd.Run(); err != nil {
		f.Close()
		return "", err
	}
	f.Close()

	tarf, err := ioutil.TempFile(tmpdir, "postgres-archive-")
	if err != nil {
		return "", err
	}
	defer tarf.Close()
	if err := createTarArchive(tarf, []string{f.Name()}); err != nil {
		return "", err
	}
	if err := os.Remove(f.Name()); err != nil {
		p.logger.Infof("Error while removing file %s", f.Name())
	}

	return tarf.Name(), nil
}
