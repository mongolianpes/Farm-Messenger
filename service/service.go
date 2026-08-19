package service

import (
	"context"
	"errors"
	"log/slog"

	pb "messenger/proto"
)

var (
	cantSendMessageYourselfError = errors.New("You cannot send a message to yourself")
	cantSendEmptyMessageError    = errors.New("You cannot send an empty message")
	cantChatWithYourselfError    = errors.New("You don't have a chat with yourself")
)

func (s *MessengerServer) GetChats(ctx context.Context, req *pb.GetChatsRequest) (*pb.GetChatsResponse, error) {
	return s.storage.GetUserChats(int(req.UserID))
}

func (s *MessengerServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if req.SenderID == req.ReceivedID {
		slog.Warn("Пользователь попытался отправить сообщение сам себе", "userID", req.SenderID)
		return nil, cantSendMessageYourselfError
	}

	if req.Text == "" {
		slog.Warn("Пользователь попытался отправить пустое сообщение", "userID", req.SenderID)
		return nil, cantSendEmptyMessageError
	}

	return s.storage.SaveMessage(int(req.SenderID), int(req.ReceivedID), int(req.RelatedAnnouncementID), req.Text)
}

func (s *MessengerServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	if req.UserID == req.PartnerID {
		slog.Warn("Пользователь попытался получить историю чата с самими собой", "userID", req.UserID)
		return nil, cantChatWithYourselfError
	}
	return s.storage.GetChatHistory(int(req.UserID), int(req.PartnerID), int(req.Offset), int(req.RelatedAnnouncementID))
}
