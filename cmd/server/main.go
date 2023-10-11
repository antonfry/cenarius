package main

import (
	"cenarius/internal/server"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	done := make(chan struct{})

	conf := server.NewConfig()
	s := server.NewServer(conf)

	go func() {
		sig := <-sigs
		log.Infof("OS SIGNAL: %v", sig)
		s.Shutdown()
		close(done)
	}()

	s.Start()
	<-done
}
