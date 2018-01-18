package datastore

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
)

// Etcd represents the datastore interface for the etcd db
type Etcd struct {
	endpoint string
}

// NewEtcd creates a new etcd datastore from a given endpoint
func NewEtcd(endpoint string) *Etcd {
	return &Etcd{
		endpoint: endpoint,
	}
}

// String returns a string representation of the datastore
func (e *Etcd) String() string {
	return "etcd"
}

// ExportTo exports the database contents to a file and uses the given tempdir
// required to satisfy the datastore interface
func (e *Etcd) ExportTo(tmpdir string) (string, error) {
	cfg := client.Config{
		Endpoints: []string{e.endpoint},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return "", err
	}

	kapi := client.NewKeysAPI(c)
	// set "/foo" key with "bar" value
	resp, err := kapi.Get(context.Background(), "/", &client.GetOptions{Recursive: true})
	if err != nil {
		return "", err
	}
	// Export and write output.
	m := etcdmap.Map(resp.Node)

	// TODO: add time stamp here
	f, err := ioutil.TempFile(tmpdir, "etcd-")
	if err != nil {
		return "", err
	}
	defer f.Close()

	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	_, err = f.Write(j)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
