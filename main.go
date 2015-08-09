package main

import (
	"flag"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

var (
	port           int
	storage        string
	dockerHost     string
	dockerCertPath string
	logFile        string
	tempDir        string
	storageDir     string
)

func init() {
	defaultPort, _ := strconv.Atoi(os.Getenv("PORT"))
	flag.IntVar(&port, "port", defaultPort, "port number")
	flag.StringVar(&storage, "storage", "filesystem", "storage type")
	flag.StringVar(&dockerHost, "docker-host", os.Getenv("DOCKER_HOST"), "docker host")
	flag.StringVar(&dockerCertPath, "docker-cert-path", os.Getenv("DOCKER_CERT_PATH"), "docker cert path")
	flag.StringVar(&logFile, "logfile", "", "log file path")
	flag.StringVar(&tempDir, "temp-dir", "workspace", "temporary dir")
	flag.StringVar(&storageDir, "storage-dir", "storage", "storage dir")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})

}

func main() {
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}

	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}
