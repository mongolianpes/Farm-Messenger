package service

import (
	"messenger/db"
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
