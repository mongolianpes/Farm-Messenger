package service

import (
	"context"
	"log/slog"

	pb "messenger/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const cantSendMessageYourself = "You cannot send a message to yourself"
const cantSendEmptyMessage = "You cannot send an empty message"
const cantChatWithYourself = "You don't have a chat with yourself"

func (s *MessengerServer) GetChats(ctx context.Context, req *pb.GetChatsRequest) (*pb.GetChatsResponse, error) {
	return s.storage.GetUserChats(int(req.UserID))
}

func (s *MessengerServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if req.SenderID == req.ReceivedID {
		slog.Warn("Пользователь попытался отправить сообщение сам себе", "userID", req.SenderID)
		return nil, status.Error(codes.InvalidArgument, cantSendMessageYourself)
	}

	if req.Text == "" {
		slog.Warn("Пользователь попытался отправить пустое сообщение", "userID", req.SenderID)
		return nil, status.Error(codes.InvalidArgument, cantSendEmptyMessage)
	}

	return s.storage.SaveMessage(int(req.SenderID), int(req.ReceivedID), int(req.RelatedAnnouncementID), req.Text)
}

func (s *MessengerServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	if req.UserID == req.PartnerID {
		slog.Warn("Пользователь попытался получить историю чата с самими собой", "userID", req.UserID)
		return nil, status.Error(codes.InvalidArgument, cantChatWithYourself)
	}
	return s.storage.GetChatHistory(int(req.UserID), int(req.PartnerID), int(req.Offset), int(req.RelatedAnnouncementID))
}
