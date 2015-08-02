package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	port           int
	storage        string
	dockerHost     string
	dockerCertPath string
)

func init() {
	flag.IntVar(&port, "port", 8080, "port number")
	flag.StringVar(&storage, "storage", "filesystem", "storage type")
	flag.StringVar(&dockerHost, "docker-host", os.Getenv("DOCKER_HOST"), "docker host")
	flag.StringVar(&dockerCertPath, "docker-cert-path", os.Getenv("DOCKER_CERT_PATH"), "docker cert path")
	flag.Parse()
}

func main() {
	err := Run()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
