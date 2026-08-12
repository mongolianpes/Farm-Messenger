package main

import (
	"log"
	"net"

	"messenger/messenger"
	pb "messenger/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		log.Fatalf("не удалось слушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessengerServer(grpcServer, messenger.NewMessengerServer())

	log.Println("gRPC сервер запущен на :8086")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}
