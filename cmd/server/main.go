package main

import "cenarius/internal/server"

func main() {
	conf := server.NewConfig()
	s := server.NewServer(conf)
	s.Start()
}
