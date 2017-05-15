package datastore

import (
	"archive/tar"
	"io/ioutil"
	"os"
	"path/filepath"
)

func createTarArchive(tarArchive *os.File, files []string) error {

	// Create a new tar archive.
	tw := tar.NewWriter(tarArchive)

	// Add some files to the archive.
	for _, filename := range files {
		body, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: filepath.Base(filename),
			Mode: 0600,
			Size: int64(len(body)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(body); err != nil {
			return err
		}
	}
	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}
