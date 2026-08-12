package messenger

import (
	"context"

	pb "messenger/proto"
)

func (s *MessengerServer) GetChats(ctx context.Context, req *pb.GetChatsRequest) (*pb.GetChatsResponse, error) {
	return getChatsFromDB(s.db, int(req.UserID))
}

func (s *MessengerServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return saveMessageInDB(s.db, int(req.SenderID), int(req.ReceivedID), int(req.RelatedAnnouncementID), req.Text)
}

func (s *MessengerServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	return getChatHistoryFromDB(s.db, int(req.UserID), int(req.PartnerID), int(req.Offset)*limitMessagesToShow, int(req.RelatedAnnouncementID))
}
