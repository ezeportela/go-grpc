package main

import (
	"log"
	"net"

	"github.com/ezeportela/go-grpc/database"
	"github.com/ezeportela/go-grpc/server"
	"github.com/ezeportela/go-grpc/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo := database.NewPostgresRepository("postgres://postgres:admin1234@localhost:5432/mydb?sslmode=disable")

	server := server.NewStudentServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	studentpb.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
