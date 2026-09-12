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
	resp, err := s.storage.GetUserChats(ctx, int(req.UserID))
	if err != nil {
		slog.Warn("Пользователь попытался получить свои чаты, но возникла ошибка", "error", err)
		return nil, err
	}

	return resp, nil
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

	resp, err := s.storage.SaveMessage(ctx, int(req.SenderID), int(req.ReceivedID), int(req.RelatedAnnouncementID), req.Text)
	if err != nil {
		slog.Warn("Пользователь попытался отправить сообщение, но возникла ошибка", "error", err)
		return nil, err
	}

	slog.Info("Успешная отправка сообщения", "senderID", req.SenderID, "receivedID", req.ReceivedID)

	return resp, nil
}

func (s *MessengerServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	if req.UserID == req.PartnerID {
		slog.Warn("Пользователь попытался получить историю чата с самими собой", "userID", req.UserID)
		return nil, cantChatWithYourselfError
	}

	resp, err := s.storage.GetChatHistory(ctx, int(req.UserID), int(req.PartnerID), int(req.Offset), int(req.RelatedAnnouncementID))
	if err != nil {
		slog.Warn("Пользователь попытался получить чат, но возникла ошибка", "error", err)
		return nil, err
	}

	slog.Info("Успешное получении истори чата", "userID", req.UserID, "partnerID", req.PartnerID)

	return resp, nil
}
