package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	port           int
	storage        string
	dockerHost     string
	dockerCertPath string
)

func init() {
	defaultPort, _ := strconv.Atoi(os.Getenv("PORT"))
	flag.IntVar(&port, "port", defaultPort, "port number")
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
