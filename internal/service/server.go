package service

import (
	"messenger/internal/db"
	pb "messenger/proto"
)

type MessengerServer struct {
	pb.UnimplementedMessengerServer
	storage db.ChatsStorage
}

func NewMessengerServer(storage db.ChatsStorage) *MessengerServer {
	return &MessengerServer{
		storage: storage,
	}
}
