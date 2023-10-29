package main

import (
	"cenarius/internal/agent"
	"cenarius/internal/server"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type flags struct {
	mode string
	conf string
}

type cenariusWorker interface {
	Start() error
	Shutdown()
}

var (
	buildVersion string
	buildDate    string
	buildCommit  string
	flagsData    flags
	worker       cenariusWorker
)

func init() {
	flag.StringVar(&flagsData.mode, "m", "", "server or agent")
	flag.StringVar(&flagsData.conf, "config", "", "Path to config")
	flag.Parse()
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	switch flagsData.mode {
	case "server":
		conf := server.NewConfig()
		worker = server.NewServer(conf)
	case "agent":
		conf := agent.NewConfig()
		worker = agent.NewAgent(conf)
	default:
		flag.Usage()
		log.Fatalf("Unknown mode %v", flagsData.mode)
	}

	go func() {
		sig := <-sigs
		log.Infof("OS SIGNAL: %v", sig)
		worker.Shutdown()
	}()
	log.Infof("Build version: %v", buildVersion)
	log.Infof("Build date: %v", buildDate)
	log.Infof("Build commit: %v", buildCommit)
	if err := worker.Start(); err != nil {
		log.Fatal(err)
	}
}
