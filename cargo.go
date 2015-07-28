package cargo

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

func (c cargo) build() error {

	var err error
	workspace, _ := ioutil.TempDir("workspace", "")
	fmt.Printf("workspace: %s\n", workspace)

	storage := newStorage("file", c.params)

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
	err = storage.save(filepath.Join(workspace, "cgotest", "app.tar.gz"))
	if err != nil {
		return buildError{err}
	}
	return nil
}
