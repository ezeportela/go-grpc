package main

import (
	"log"

	"github.com/ezeportela/go-grpc/server"
	"github.com/ezeportela/go-grpc/testpb"
	"google.golang.org/grpc"
)

func main() {
	listener, repo := server.SetupServer()

	srv := server.NewTestServer(repo)

	s := server.NewGrpcServer(func(s *grpc.Server) {
		testpb.RegisterTestServiceServer(s, srv)
	})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
