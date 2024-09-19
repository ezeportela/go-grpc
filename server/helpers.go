package server

import (
	"log"
	"net"
	"os"

	"github.com/ezeportela/go-grpc/database"
	"github.com/ezeportela/go-grpc/repositories"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetupServer() (net.Listener, repositories.Repository) {
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo := database.NewPostgresRepository(dbUrl)

	return listener, repo
}

func NewGrpcServer(register func(*grpc.Server)) *grpc.Server {
	s := grpc.NewServer()

	register(s)

	reflection.Register(s)

	return s
}
