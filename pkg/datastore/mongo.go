package datastore

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

const (
	mongodump = "mongodump"
)

func mkMongodumpCmd(host string, port int, user, password string, archFile string) *exec.Cmd {
	args := []string{
		"--host", host,
		"--port", fmt.Sprintf("%d", port),
		fmt.Sprintf("--archive=%s", archFile),
		"--gzip",
	}
	if user != "" {
		args = append(args, "--user", user)
	}
	if password != "" {
		args = append(args, "--password", password)
	}

	return exec.Command(mongodump, args...)
}

// MongoDB represents the datastore interface for the mongodb database
type MongoDB struct {
	host     string
	port     int
	user     string
	password string
}

// NewMongoDB creates a new MongoDB instance from the given settings
func NewMongoDB(host string, port int, user, password string) *MongoDB {
	return &MongoDB{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}
}

// String returns a string representation of the datastore
func (m *MongoDB) String() string {
	return "mongodb"
}

// ExportTo exports the database contents to a file and uses the given tempdir
// required to satisfy the datastore interface
func (m *MongoDB) ExportTo(tmpdir string) (string, error) {

	f, err := ioutil.TempFile(tmpdir, "mongo-")
	if err != nil {
		return "", err
	}
	tmpfile := f.Name()
	defer f.Close()

	cmd := mkMongodumpCmd(m.host, m.port, m.user, m.password, tmpfile)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	return tmpfile, nil
}
