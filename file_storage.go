package main

import (
	"os"
	"path"
	"path/filepath"
)

func newFileStorage(params *params) storager {
	return &fileStorage{params}
}

type fileStorage struct {
	params *params
}

func (f fileStorage) isExist() bool {
	if _, err := os.Stat(f.pathByParams()); err != nil {
		return false
	}
	return true
}

func (f fileStorage) save(from string) error {
	if err := os.MkdirAll(f.pathByParams(), 0755); err != nil {
		return err
	}
	return os.Rename(from, filepath.Join(f.pathByParams(), path.Base(from)))
}

func (f fileStorage) pathByParams() string {
	return filepath.Join(
		"storage",
		f.params.remote,
		f.params.owner(),
		f.params.repo,
		f.params.goos,
		f.params.goarch,
		f.params.version,
	)
}

func (f fileStorage) get(file string) (string, error) {
	return filepath.Join(f.pathByParams(), file), nil
}
