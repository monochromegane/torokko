package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type cargo struct {
	params *params
}

func newCargo(params map[string]string) *cargo {
	return &cargo{newParams(params)}
}

func (c cargo) store(queue chan *params) error {
	storage := newStorage(c.params)

	// exist?
	if storage.isExist() {
		return aleadyExistsError{}
	}

	// store in build queue
	queue <- c.params

	return nil
}

func (c cargo) build() error {

	var err error
	workspace, _ := ioutil.TempDir("workspace", "")
	fmt.Printf("workspace: %s\n", workspace)

	storage := newStorage(c.params)

	// exist?
	if storage.isExist() {
		return aleadyExistsError{}
	}

	// clone
	repo := newRepository(c.params, workspace)
	err = repo.clone("https")
	if err != nil {
		return buildError{err}
	}

	// build
	builder := newBuilder(c.params)
	err = builder.build(workspace)
	if err != nil {
		return buildError{err}
	}

	// diff archive
	err = repo.diffArchive("app", "tar.gz")
	if err != nil {
		return buildError{err}
	}

	// save
	err = storage.save(filepath.Join(workspace, c.params.repo, "app.tar.gz"))
	if err != nil {
		return buildError{err}
	}
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
