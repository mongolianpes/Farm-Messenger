package main

import (
	"log/slog"
	"net"

	"messenger/internal/db"
	"messenger/internal/service"
	pb "messenger/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	storage, err := db.NewChatsStorage()
	if err != nil {
		panic(err)
	}

	pb.RegisterMessengerServer(grpcServer, service.NewMessengerServer(storage))

	slog.Info("Сервер запущен")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
