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

	fmt.Println(args)
	return exec.Command(mongodump, args...)
}

type MongoDB struct {
	host     string
	port     int
	user     string
	password string
}

func NewMongoDB(host string, port int, user, password string) *MongoDB {
	return &MongoDB{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}
}

func (m *MongoDB) ExportTo(tmpdir string) (string, error) {

	f, err := ioutil.TempFile(tmpdir, "mongo-")
	if err != nil {
		return "", err
	}
	tmpfile := f.Name()
	defer f.Close()

	cmd := mkMongodumpCmd(m.host, m.port, m.user, m.password, tmpfile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return "", err
	}

	return tmpfile, nil
}