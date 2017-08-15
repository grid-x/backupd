package datastore

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	influxd = "influxd"
)

func mkInfluxdCmd(database, host, tmpdir string) *exec.Cmd {
	return exec.Command(influxd, "backup", "-host", host,
		"-database", database, tmpdir)
}

type Influx struct {
	endpoint string
	database string
}

func NewInflux(endpoint, database string) *Influx {
	return &Influx{
		endpoint: endpoint,
		database: database,
	}
}

func (i *Influx) ExportTo(tmpdir string) (string, error) {
	cmd := mkInfluxdCmd(i.database, i.endpoint, tmpdir)

	//TODO: do not ignore output, but log it in case of error
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
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
		os.Remove(f)
	}

	return archive.Name(), nil
}
