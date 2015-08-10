package main

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	dockerclient "github.com/fsouza/go-dockerclient"
)

func newDocker(opt *containerOption, logger *log.Entry) *docker {
	return &docker{opt: opt, client: client(), logger: logger}
}

type docker struct {
	opt    *containerOption
	client *dockerclient.Client
	logger *log.Entry
}

type containerOption struct {
	image      string
	env        []string
	cmd        []string
	volumes    []string
	workingDir string
}

func (c containerOption) toLog() map[string]interface{} {
	return map[string]interface{}{
		"image":      c.image,
		"env":        c.env,
		"cmd":        c.cmd,
		"volumes":    c.volumes,
		"workingDir": c.workingDir,
	}
}

func (d docker) run() (int, error) {

	d.logger.Info("creating container...")
	con, err := d.createContainer()
	if err != nil {
		d.logger.Warnf("creating container error: %v", err)
		return 0, err
	}

	log := d.logger.WithField("container", con.ID)

	log.Info("starting container...")
	err = d.startContainer(con.ID)
	if err != nil {
		log.Warnf("starting container error: %v", err)
		return 0, err
	}

	log.Info("waiting container...")
	status, err := d.waitContainer(con.ID)
	if err != nil {
		log.Warnf("waiting container error: %v", err)
		return 0, err
	}

	err = d.loggingContainerLog(con.ID)
	if err != nil {
		log.Warnf("logging container error: %v", err)
		return 0, err
	}

	log.Info("removing container...")
	err = d.removeContainer(con.ID)
	if err != nil {
		log.Warnf("removing container error: %v", err)
		return 0, err
	}
	return status, nil
}

func (d docker) createContainer() (*dockerclient.Container, error) {
	volumes := map[string]struct{}{}
	for _, v := range d.opt.volumes {
		volumes[v] = struct{}{}
	}
	return d.client.CreateContainer(dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image:      d.opt.image,
			Env:        d.opt.env,
			Cmd:        d.opt.cmd,
			Volumes:    volumes,
			WorkingDir: d.opt.workingDir,
		},
		HostConfig: &dockerclient.HostConfig{
			Binds: d.opt.volumes,
		},
	})
}

func (d docker) startContainer(id string) error {
	return d.client.StartContainer(id, nil)
}

func (d docker) waitContainer(id string) (int, error) {
	return d.client.WaitContainer(id)
}

func (d docker) removeContainer(id string) error {
	return d.client.RemoveContainer(
		dockerclient.RemoveContainerOptions{ID: id},
	)
}

func (d docker) inspectContainer(id string) (*dockerclient.Container, error) {
	return d.client.InspectContainer(id)
}

func (d docker) loggingContainerLog(id string) error {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	err := d.client.Logs(dockerclient.LogsOptions{
		Container:    id,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stdout:       true,
		Stderr:       true,
		Follow:       false,
	})
	if err != nil {
		return err
	}
	d.logger.WithFields(log.Fields{
		"container": id,
		"command":   "/containers/(id)/logs",
		"args":      []string{"GET", fmt.Sprintf("/containers/%s/logs", id)},
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
	}).Info("logging container...")
	return nil
}

func client() *dockerclient.Client {
	if dockerCertPath != "" {
		return tlsClient()
	} else {
		client, _ := dockerclient.NewClient(dockerHost)
		return client
	}
}

func tlsClient() *dockerclient.Client {
	ca := fmt.Sprintf("%s/ca.pem", dockerCertPath)
	cert := fmt.Sprintf("%s/cert.pem", dockerCertPath)
	key := fmt.Sprintf("%s/key.pem", dockerCertPath)
	client, _ := dockerclient.NewTLSClient(dockerHost, cert, key, ca)
	return client
}
