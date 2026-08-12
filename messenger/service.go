package messenger

import (
	"context"
	"time"

	pb "messenger/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const limitMessagesToShow = 5

func (s *MessengerServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if _, err := s.db.Exec("INSERT INTO messages (sender_id, received_id, related_announcement_id, message) VALUES ($1, $2, $3)", req.SenderID, req.ReceivedID, req.RelatedAnnouncementID, req.Text); err != nil {
		return &pb.SendMessageResponse{
			Error: err.Error(),
		}, err
	}

	return &pb.SendMessageResponse{}, nil
}

func (s *MessengerServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	query := "SELECT create_at, message, sender_id FROM messages WHERE  AND related_announcement_id = $5 ((sender_id = $1 AND received_id = $2) OR (sender_id = $2 AND received_id = $1)) ORDER BY create_at DESC OFFSET $3 LIMIT $4"
	rows, err := s.db.Query(query, req.UserID, req.PartnerID, req.Offset*limitMessagesToShow, limitMessagesToShow, req.RelatedAnnouncementID)
	if err != nil {
		return &pb.GetChatHistoryResponse{}, err
	}

	var messages []*pb.GetChatHistoryMessage
	var myMessage bool
	var createAt time.Time
	var text string
	var senderID int
	for rows.Next() {
		if err := rows.Scan(&createAt, &text, &senderID); err != nil {
			return &pb.GetChatHistoryResponse{}, err
		}

		if senderID == int(req.UserID) {
			myMessage = true
		} else {
			myMessage = false
		}

		messages = append(messages, &pb.GetChatHistoryMessage{
			CreateAt:  timestamppb.New(createAt),
			Text:      text,
			MyMessage: myMessage,
		})
	}

	return &pb.GetChatHistoryResponse{
		Messages: messages,
	}, nil
}
