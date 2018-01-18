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

type Postgres struct {
	logger log.FieldLogger
	url    string
}

func NewPostgres(url string) *Postgres {
	return &Postgres{
		logger: log.New().WithFields(log.Fields{
			"component": "datastore",
			"datastore": "postgres",
		}),
		url: url,
	}
}

func (p *Postgres) String() string {
	return "postgresql"
}

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
