package db

import (
	"database/sql"
	pb "messenger/proto"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const limitMessagesToShow = 5
const queryGetChats = `
SELECT DISTINCT 
    CASE 
        WHEN m.sender_id = $1 THEN m.received_id 
        ELSE m.sender_id 
    END AS partner_id,
    u.name,
    u.avatar_path,
    a.title,
	m.related_announcement_id
FROM messages m
JOIN users u 
    ON u.user_id = CASE 
        WHEN m.sender_id = $1 THEN m.received_id 
        ELSE m.sender_id 
    END
JOIN announcements a 
    ON a.announcement_id = m.related_announcement_id
WHERE m.sender_id = $1 OR m.received_id = $1;`

type PostgresStorage struct {
	db *sql.DB
}

type ChatsStorage interface {
	GetUserChats(userID int) (*pb.GetChatsResponse, error)
	SaveMessage(senderID, receivedID, relatedAnnouncementID int, text string) (*pb.SendMessageResponse, error)
	GetChatHistory(userID, partnerID, offset, relatedAnnouncementID int) (*pb.GetChatHistoryResponse, error)
}

func (s *PostgresStorage) GetUserChats(userID int) (*pb.GetChatsResponse, error) {
	rows, err := s.db.Query(queryGetChats, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*pb.ChatInfo
	var partnerID int
	var partnerName string
	var partnerAvatarPath string
	var relatedAnnouncementTitle string
	var relatedAnnouncementID int
	for rows.Next() {
		if err := rows.Scan(&partnerID, &partnerName, &partnerAvatarPath, &relatedAnnouncementTitle, &relatedAnnouncementID); err != nil {
			return nil, err
		}

		chats = append(chats, &pb.ChatInfo{
			PartnerID:                int32(partnerID),
			PartnerName:              partnerName,
			RelatedAnnouncementTitle: relatedAnnouncementTitle,
			PartnerAvatarPath:        partnerAvatarPath,
			RelatedAnnouncementID:    int32(relatedAnnouncementID),
		})
	}

	return &pb.GetChatsResponse{
		Chats: chats,
	}, nil
}

func (s *PostgresStorage) SaveMessage(senderID, receivedID, relatedAnnouncementID int, text string) (*pb.SendMessageResponse, error) {
	if _, err := s.db.Exec("INSERT INTO messages (sender_id, received_id, related_announcement_id, message) VALUES ($1, $2, $3, $4)", senderID, receivedID, relatedAnnouncementID, text); err != nil {
		return nil, err
	}

	return &pb.SendMessageResponse{}, nil
}

func (s *PostgresStorage) GetChatHistory(userID, partnerID, offset, relatedAnnouncementID int) (*pb.GetChatHistoryResponse, error) {
	query := "SELECT create_at, message, sender_id FROM messages WHERE ((sender_id = $1 AND received_id = $2) OR (sender_id = $2 AND received_id = $1)) AND related_announcement_id = $5 ORDER BY create_at DESC OFFSET $3 LIMIT $4"
	rows, err := s.db.Query(query, userID, partnerID, offset*limitMessagesToShow, limitMessagesToShow, relatedAnnouncementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*pb.GetChatHistoryMessage
	var myMessage bool
	var createAt time.Time
	var text string
	var senderID int
	for rows.Next() {
		if err := rows.Scan(&createAt, &text, &senderID); err != nil {
			return nil, err
		}

		if senderID == int(userID) {
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
