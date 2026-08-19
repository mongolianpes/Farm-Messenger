package main

import (
	"log/slog"
	"net"

	"messenger/db"
	pb "messenger/proto"
	"messenger/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		slog.Error("не удалось слушать порт", "error", err)
	}

	grpcServer := grpc.NewServer()

	storage, err := db.NewChatsStorage()
	if err != nil {
		panic(err)
	}

	pb.RegisterMessengerServer(grpcServer, service.NewMessengerServer(storage))

	slog.Info("Сервер запущен")
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Ошибка при запуске сервера", "error", err)
	}
}
