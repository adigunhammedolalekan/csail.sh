package main

import (
	"github.com/saas/hostgolang/pkg/server"
	"log"
)

func main() {
	s, err := server.NewServer(":4005")
	if err != nil {
		 log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
