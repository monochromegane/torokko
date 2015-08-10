package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type builder struct {
	params *params
	logger *log.Entry
}

func newBuilder(params *params, logger *log.Entry) *builder {
	return &builder{params: params, logger: logger}
}

func (b builder) build(dir string) error {
	docker := newDocker(&containerOption{
		image: "golang:1.4.2-cross",
		env: []string{
			"GOOS=" + b.params.goos,
			"GOARCH=" + b.params.goarch,
		},
		cmd:        []string{"make"},
		volumes:    []string{b.volumeFrom(dir) + ":" + b.volumeTo()},
		workingDir: b.volumeTo(),
	}, b.logger)
	b.addMakefileUnlessExists(dir)
	_, err := docker.run()
	if err != nil {
		return err
	}
	return nil
}

func (b builder) volumeFrom(dir string) string {
	path, _ := filepath.Abs(dir)
	return filepath.Join(path, b.params.repo)
}

func (b builder) volumeTo() string {
	return filepath.Join(
		"/go/src",
		b.params.remote,
		b.params.owner(),
		b.params.repo,
	)
}

func (b builder) addMakefileUnlessExists(dir string) {
	makefile := filepath.Join(b.volumeFrom(dir), "Makefile")
	if _, err := os.Stat(makefile); err != nil {
		b.logger.Info("adding a default Makefile...")
		f, _ := os.Create(makefile)
		f.WriteString(
			`build:
	go get -d ./...
	go build
`)
	}
}
