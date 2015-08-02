package main

import (
	"fmt"
	"os"

	dockerclient "github.com/fsouza/go-dockerclient"
)

func newDocker(opt *containerOption) *docker {
	return &docker{opt: opt, client: client()}
}

type docker struct {
	opt    *containerOption
	client *dockerclient.Client
}

type containerOption struct {
	image      string
	env        []string
	cmd        []string
	volumes    []string
	workingDir string
}

func (d docker) run() (int, error) {
	con, err := d.createContainer()
	if err != nil {
		return 0, err
	}

	err = d.startContainer(con.ID)
	if err != nil {
		return 0, err
	}

	status, err := d.waitContainer(con.ID)
	if err != nil {
		return 0, err
	}

	err = d.removeContainer(con.ID)
	if err != nil {
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

func client() *dockerclient.Client {
	endpoint := os.Getenv("DOCKER_HOST")
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, _ := dockerclient.NewTLSClient(endpoint, cert, key, ca)
	return client
}
