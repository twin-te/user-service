package main

import (
	"log"
	"net"

	"github.com/twin-te/user-service/repository"
	"github.com/twin-te/user-service/server"
	"github.com/twin-te/user-service/usecase"
)

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	r, close, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	u := usecase.NewUseCase(r)
	s := server.NewServer(u)

	log.Printf("server listening at %v", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
