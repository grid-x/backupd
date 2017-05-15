package datastore

import (
	"os/exec"
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

func (i *Influx) ExportTo(tmpdir string) error {
	cmd := mkInfluxdCmd(i.database, i.endpoint, tmpdir)
	//TODO: do not ignore output, but log it in case of error
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
