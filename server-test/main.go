package main

import (
	"log"
	"net"
	"os"

	"github.com/ezeportela/go-grpc/database"
	"github.com/ezeportela/go-grpc/server"
	"github.com/ezeportela/go-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo := database.NewPostgresRepository(dbUrl)

	server := server.NewTestServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	testpb.RegisterTestServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
