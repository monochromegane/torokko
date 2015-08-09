package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type cargo struct {
	params *params
	logger *log.Entry
}

func newCargo(params *params) *cargo {
	log := log.WithFields(log.Fields{
		"params": params.params,
	})
	return &cargo{
		params: params,
		logger: log,
	}
}

func (c cargo) store(queue chan *params) error {
	storage := newStorage(c.params)

	// exist?
	if storage.isExist() {
		return aleadyExistsError{}
	}

	// store in build queue
	queue <- c.params
	c.logger.Info("stored")

	return nil
}

func (c cargo) build() error {

	var err error
	workspace, err := ioutil.TempDir("workspace", "")
	if err != nil {
		return err
	}
	c.logger.Infof("build start on %s", workspace)

	storage := newStorage(c.params)

	// exist?
	if storage.isExist() {
		c.logger.Warn("already exists")
		return aleadyExistsError{}
	}

	// clone
	repo := newRepository(c.params, workspace)
	err = repo.clone("https")
	if err != nil {
		c.logger.Warnf("git clone error: %v", err)
		return buildError{err}
	}

	// build
	builder := newBuilder(c.params)
	err = builder.build(workspace)
	if err != nil {
		c.logger.Warnf("build error: %v", err)
		return buildError{err}
	}

	// diff archive
	err = repo.diffArchive("app", "tar.gz")
	if err != nil {
		c.logger.Warnf("diff archive error: %v", err)
		return buildError{err}
	}

	// save
	err = storage.save(filepath.Join(workspace, c.params.repo, "app.tar.gz"))
	if err != nil {
		c.logger.Warnf("save error: %v", err)
		return buildError{err}
	}
	c.logger.Info("build success")
	return nil
}

func (c cargo) isExist() bool {
	return newStorage(c.params).isExist()
}

func (c cargo) get() (string, error) {
	return newStorage(c.params).get("app.tar.gz")
}

func (c cargo) downloadFileName() string {
	return fmt.Sprintf(
		"%s_%s_%s.tar.gz",
		c.params.repo,
		c.params.goos,
		c.params.goarch,
	)
}
