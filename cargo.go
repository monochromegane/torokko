package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

type cargo struct {
	params *params
	logger *log.Entry
}

func newCargo(params *params) *cargo {
	return &cargo{params: params}
}

func (c *cargo) store(queue chan *params) (string, error) {
	storage := newStorage(c.params)

	// exist?
	if storage.isExist() {
		return "", aleadyExistsError{}
	}

	// store in build queue
	c.params.buildId = c.genBuildId()
	queue <- c.params

	f, err := c.openBuildLog()
	if err != nil {
		return "", err
	}
	defer f.Close()
	c.setBuildLogger(f)

	c.logger.WithFields(
		log.Fields{"params": c.params.params},
	).Info("Your build joined a queue.")

	return c.params.buildId, nil
}

func (c cargo) build() error {

	f, err := c.openBuildLog()
	if err != nil {
		return err
	}
	defer f.Close()
	c.setBuildLogger(f)

	workspace, err := ioutil.TempDir(tempDir, "")
	if err != nil {
		return err
	}
	c.logger.WithField("workspace", filepath.Base(workspace)).Info("Your build started.")

	storage := newStorage(c.params)

	// exist?
	c.logger.Info("checking binary...")
	if storage.isExist() {
		c.logger.Info("The Binary already exists.")
		return aleadyExistsError{}
	}

	// clone
	c.logger.Info("cloning repository...")
	repo := newRepository(c.params, workspace, c.logger)
	err = repo.clone("https")
	if err != nil {
		return buildError{err}
	}

	// build
	c.logger.Info("building binary...")
	builder := newBuilder(c.params, c.logger)
	err = builder.build(workspace)
	if err != nil {
		return buildError{err}
	}

	// diff archive
	c.logger.Info("archiving binary...")
	err = repo.diffArchive("app", "tar.gz")
	if err != nil {
		return buildError{err}
	}

	// save
	c.logger.Info("saving binary...")
	err = storage.save(filepath.Join(workspace, c.params.repo, "app.tar.gz"))
	if err != nil {
		return buildError{err}
	}

	c.logger.Info("Your build is successful.")
	return nil
}

func (c *cargo) setBuildLogger(f *os.File) {
	var logger = log.New()
	logger.Formatter = &log.JSONFormatter{}
	logger.Out = f
	c.logger = logger.WithFields(log.Fields{"build_id": c.params.buildId})
}

func (c cargo) openBuildLog() (*os.File, error) {
	return os.OpenFile(filepath.Join(logDir, c.params.buildId), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func (c cargo) isExist() bool {
	return newStorage(c.params).isExist()
}

func (c cargo) isAuthorized() bool {
	repo := newRepository(c.params, "", nil)
	err := repo.listRemote()
	if err != nil {
		return false
	}
	return true
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

func (c cargo) genBuildId() string {
	source := fmt.Sprintf("%s/%s/%s/%s/%s/%s+%s",
		c.params.remote,
		c.params.owner(),
		c.params.repo,
		c.params.goos,
		c.params.goarch,
		c.params.version,
		time.Now().Format("20060102150405"),
	)
	hasher := md5.New()
	hasher.Write([]byte(source))
	return hex.EncodeToString(hasher.Sum(nil))
}
