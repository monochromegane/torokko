package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type logFile struct {
	id string
}

func (l logFile) isExist() bool {
	if _, err := os.Stat(l.path()); err != nil {
		return false
	}
	return true
}

func (l logFile) readAll() ([]byte, error) {
	return ioutil.ReadFile(l.path())
}

func (l logFile) path() string {
	return filepath.Join(logDir, l.id)
}
