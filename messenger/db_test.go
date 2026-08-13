package messenger

import (
	"database/sql"
	"testing"
	"time"
)

func deleteTestMessage(db *sql.DB, messageText string) error {
	if _, err := db.Exec("DELETE FROM messages WHERE message = $1", messageText); err != nil {
		return err
	}

	return nil
}

type updMessagesInfo struct {
	usersIDs         []int
	announcementsIDs []int
}

func updateTestUsers(db *sql.DB) (updMessagesInfo, error) {
	updMessages := updMessagesInfo{}
	if _, err := db.Exec("DELETE FROM users"); err != nil {
		return updMessages, err
	}

	rowsUserIDs, err := db.Query(`INSERT INTO users (login, name, password) VALUES
		(1, '1', '1234567890'),
		(2, '2', '1234567890'),
		(3, '3', '1234567890')
		RETURNING user_id`)
	if err != nil {
		return updMessages, err
	}

	var userID1 int
	var userID2 int
	var userID3 int
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID1)
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID2)
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID3)
	updMessages.usersIDs = append(updMessages.usersIDs, userID1)
	updMessages.usersIDs = append(updMessages.usersIDs, userID2)
	updMessages.usersIDs = append(updMessages.usersIDs, userID3)

	return updMessages, err
}

func updateTestAnnouncements(db *sql.DB, updMessages updMessagesInfo) (updMessagesInfo, error) {
	rowsAnnouncementIDs, err := db.Query(`INSERT INTO announcements (title, description, category, announcement_author_id) VALUES
	('1', '1', 'Участки', $1),
	('2', '2', 'Участки', $2),
	('22', '22', 'Участки', $2),
	('222', '222', 'Участки', $2)
	RETURNING announcement_id`, updMessages.usersIDs[0], updMessages.usersIDs[1])
	if err != nil {
		return updMessages, err
	}

	var announcementID1 int
	var announcementID2 int
	var announcementID3 int
	var announcementID4 int
	rowsAnnouncementIDs.Next()
	rowsAnnouncementIDs.Scan(&announcementID1)
	rowsAnnouncementIDs.Next()
	rowsAnnouncementIDs.Scan(&announcementID2)
	rowsAnnouncementIDs.Next()
	rowsAnnouncementIDs.Scan(&announcementID3)
	rowsAnnouncementIDs.Next()
	rowsAnnouncementIDs.Scan(&announcementID4)
	updMessages.announcementsIDs = append(updMessages.announcementsIDs, announcementID1)
	updMessages.announcementsIDs = append(updMessages.announcementsIDs, announcementID2)
	updMessages.announcementsIDs = append(updMessages.announcementsIDs, announcementID3)
	updMessages.announcementsIDs = append(updMessages.announcementsIDs, announcementID4)

	return updMessages, nil
}

func updateTestMessages(db *sql.DB) (updMessagesInfo, error) {
	db.Exec("DELETE FROM messages")

	updMessages, err := updateTestUsers(db)
	if err != nil {
		return updMessages, err
	}

	updMessages, err = updateTestAnnouncements(db, updMessages)
	if err != nil {
		return updMessages, err
	}

	messages := []struct {
		senderID, receivedID, announcementID int
		text                                 string
	}{
		{updMessages.usersIDs[0], updMessages.usersIDs[1], updMessages.announcementsIDs[0], "Привет по объявлению 10"},
		{updMessages.usersIDs[1], updMessages.usersIDs[0], updMessages.announcementsIDs[0], "Ответ по объявлению 10"},
		{updMessages.usersIDs[0], updMessages.usersIDs[1], updMessages.announcementsIDs[0], "Привет по объявлению 10 (2)"},
		{updMessages.usersIDs[0], updMessages.usersIDs[1], updMessages.announcementsIDs[0], "Привет по объявлению 10 (3)"},
		{updMessages.usersIDs[1], updMessages.usersIDs[0], updMessages.announcementsIDs[0], "Ответ по объявлению 10 (2)"},
		{updMessages.usersIDs[1], updMessages.usersIDs[0], updMessages.announcementsIDs[0], "Ответ по объявлению 10 (3)"},
		{updMessages.usersIDs[0], updMessages.usersIDs[1], updMessages.announcementsIDs[1], "Привет по объявлению 20"},
		{updMessages.usersIDs[0], updMessages.usersIDs[1], updMessages.announcementsIDs[2], "Прив по 3 для 2"},
		{updMessages.usersIDs[0], updMessages.usersIDs[2], updMessages.announcementsIDs[2], "Прив по 3 для 3"},
		{updMessages.usersIDs[0], updMessages.usersIDs[2], updMessages.announcementsIDs[3], "Прив по 4 для 3"},
		{updMessages.usersIDs[0], updMessages.usersIDs[2], updMessages.announcementsIDs[0], "Сообщение пользователю 3 по объявлению 10"},
		{updMessages.usersIDs[2], updMessages.usersIDs[0], updMessages.announcementsIDs[0], "Ответ пользователю 1 по объявлению 10"},
	}

	for _, m := range messages {
		if _, err := db.Exec(`INSERT INTO messages (sender_id, received_id, related_announcement_id, message)
							VALUES ($1, $2, $3, $4)`,
			m.senderID, m.receivedID, m.announcementID, m.text); err != nil {
			return updMessages, err
		}
		time.Sleep(50 * time.Millisecond) // задержка между вставками
	}

	return updMessages, nil
}

func TestGetChatsFromDB(t *testing.T) {
	server := NewMessengerServer()
	db := server.db

	updMessages, err := updateTestMessages(db)
	if err != nil {
		t.Error(err)
	}

	userID := updMessages.usersIDs[0]
	response, err := getChatsFromDB(db, userID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Chats) != 3 {
		t.Errorf("Количество чатов не соответствуюет нужному. UserID %v", userID)
	}

	userID = updMessages.usersIDs[1]
	response, err = getChatsFromDB(db, userID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Chats) != 2 {
		t.Errorf("Количество чатов не соответствуюет нужному. UserID %v", userID)
	}

	userID = updMessages.usersIDs[2]
	response, err = getChatsFromDB(db, userID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Chats) != 1 {
		t.Errorf("Количество чатов не соответствуюет нужному. UserID %v", userID)
	}
}

func TestSaveMessageInDB(t *testing.T) {
	server := NewMessengerServer()
	db := server.db

	updMessages, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	updMessages, err = updateTestAnnouncements(db, updMessages)
	if err != nil {
		t.Error(err)
	}

	senderID := updMessages.usersIDs[0]
	receivedID := updMessages.usersIDs[1]
	relatedAnnounID := updMessages.announcementsIDs[0]
	messageText := "Здарова, это тестовое сообщение"
	_, err = saveMessageInDB(db, senderID, receivedID, relatedAnnounID, messageText)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := deleteTestMessage(db, messageText); err != nil {
			t.Error(err)
		}
	}()

	var querySenderID int
	var queryReceivedID int
	var queryRelatedAnnounID int
	var query = "SELECT sender_id, received_id, related_announcement_id FROM messages WHERE message = $1"
	rows, err := db.Query(query, messageText)
	if err != nil {
		t.Error(err)
	}

	rows.Next()
	if rows.Scan(&querySenderID, &queryReceivedID, &queryRelatedAnnounID); err != nil {
		t.Error(err)
	}
	if rows.Next() {
		t.Error("Получилось отправить сообщение самому себе или пустое сообщение")
	}

	if senderID != querySenderID || receivedID != queryReceivedID || relatedAnnounID != queryRelatedAnnounID {
		t.Error("Одно из значение не равно тому, что должно быть")
	}
}

func TestGetChatHistoryFromDB(t *testing.T) {
	server := NewMessengerServer()
	db := server.db

	updMessages, err := updateTestMessages(db)
	if err != nil {
		t.Error(err)
	}

	userID := updMessages.usersIDs[0]
	partnerID := updMessages.usersIDs[1]
	offset := 0
	relatedAnnouncementID := updMessages.announcementsIDs[0]
	response, err := getChatHistoryFromDB(db, userID, partnerID, offset, relatedAnnouncementID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Messages) != 5 {
		t.Errorf("Количество полученных сообщений не соответствует нужному. Получено сообщений: %v, offset %v", len(response.Messages), offset)
	}

	offset++
	response, err = getChatHistoryFromDB(db, userID, partnerID, offset, relatedAnnouncementID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Messages) != 1 {
		t.Errorf("Количество полученных сообщений не соответствует нужному. Получено сообщений: %v, offset %v, userID %v, partnerID %v, relatedAnnouncementID %v", len(response.Messages), offset, userID, partnerID, relatedAnnouncementID)
	}

	userID = updMessages.usersIDs[1]
	partnerID = updMessages.usersIDs[0]
	offset = 0
	relatedAnnouncementID = updMessages.announcementsIDs[1]
	response, err = getChatHistoryFromDB(db, userID, partnerID, offset, relatedAnnouncementID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Messages) != 1 {
		t.Errorf("Количество полученных сообщений не соответствует нужному. Получено сообщений: %v, offset %v, userID %v, partnerID %v, relatedAnnouncementID %v", len(response.Messages), offset, userID, partnerID, relatedAnnouncementID)
	}

	userID = updMessages.usersIDs[2]
	partnerID = updMessages.usersIDs[0]
	offset = 0
	relatedAnnouncementID = updMessages.announcementsIDs[0]
	response, err = getChatHistoryFromDB(db, userID, partnerID, offset, relatedAnnouncementID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Messages) != 2 {
		t.Errorf("Количество полученных сообщений не соответствует нужному. Получено сообщений: %v, offset %v, userID %v, partnerID %v, relatedAnnouncementID %v", len(response.Messages), offset, userID, partnerID, relatedAnnouncementID)
	}

	userID = updMessages.usersIDs[2]
	partnerID = updMessages.usersIDs[0]
	offset = 0
	relatedAnnouncementID = updMessages.announcementsIDs[1]
	response, err = getChatHistoryFromDB(db, userID, partnerID, offset, relatedAnnouncementID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Messages) != 0 {
		t.Errorf("Количество полученных сообщений не соответствует нужному. Получено сообщений: %v, offset %v, userID %v, partnerID %v, relatedAnnouncementID %v", len(response.Messages), offset, userID, partnerID, relatedAnnouncementID)
	}

}
