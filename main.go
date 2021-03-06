package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

const version = "0.0.1"

var (
	port           int
	storage        string
	dockerHost     string
	dockerCertPath string
	logDir         string
	tempDir        string
	storageDir     string
	v              bool
)

func init() {
	defaultPort, _ := strconv.Atoi(os.Getenv("PORT"))
	flag.BoolVar(&v, "version", false, "version")
	flag.IntVar(&port, "port", defaultPort, "port number")
	flag.StringVar(&storage, "storage", "filesystem", "storage type")
	flag.StringVar(&dockerHost, "docker-host", os.Getenv("DOCKER_HOST"), "docker host")
	flag.StringVar(&dockerCertPath, "docker-cert-path", os.Getenv("DOCKER_CERT_PATH"), "docker cert path")
	flag.StringVar(&logDir, "log-dir", "", "log file dir")
	flag.StringVar(&tempDir, "temp-dir", "workspace", "temporary dir")
	flag.StringVar(&storageDir, "storage-dir", "storage", "storage dir")
	flag.Parse()
}

func main() {
	if v {
		fmt.Printf("Torokko version %s\n", version)
		os.Exit(0)
	}
	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}
