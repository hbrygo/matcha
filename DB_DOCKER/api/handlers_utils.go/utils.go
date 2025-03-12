package handlers_utils

import (
	"database/sql"
	//"matcha/api/handlers"
	"matcha/api/models"
)

// check if users exist in database
func CheckUsersExist(db *sql.DB, userIDs []int) bool {
	for _, uid := range userIDs {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE uid = ?)", uid).Scan(&exists)
		if err != nil || !exists {
			return false
		}
	}
	return true
}

// create a new chatroom and return its id
func CreateNewChatroom(tx *sql.Tx, name string) (int, error) {
	result, err := tx.Exec("INSERT INTO chat_rooms (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}

	// get room id
	chatroomID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(chatroomID), nil
}

// add participants to a chatroom
func AddParticipantsToRoom(tx *sql.Tx, chatroomID int, userIDs []int) error {
	for _, uid := range userIDs {
		_, err := tx.Exec("INSERT INTO chat_participants (chatroom_id, user_uid) VALUES (?, ?)",
			chatroomID, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

// check if chatroom exists
func CheckChatroomExists(db *sql.DB, chatroomID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_rooms WHERE chatroom_id = ?)", chatroomID).Scan(&exists)
	return err == nil && exists
}

// check if user is already in chatroom
func IsUserInChatroom(db *sql.DB, chatroomID int, userID int) bool {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM chat_participants WHERE chatroom_id = ? AND user_uid = ?)",
		chatroomID, userID).Scan(&exists)
	return err == nil && exists
}

// remove user from chatroom
func RemoveUserFromChatroom(tx *sql.Tx, chatroomID int, userID int) error {
	_, err := tx.Exec(
		"DELETE FROM chat_participants WHERE chatroom_id = ? AND user_uid = ?",
		chatroomID, userID)
	return err
}

// add message to chatroom
func AddMessageToChatroom(tx *sql.Tx, chatroomID int, userID int, content string) error {
	_, err := tx.Exec(
		"INSERT INTO messages (chatroom_id, user_uid, content) VALUES (?, ?, ?)",
		chatroomID, userID, content)
	return err
}

// get messages from chatroom
func GetChatroomMessages(db *sql.DB, chatroomID int) ([]models.Message, error) {
	rows, err := db.Query(`
        SELECT m.user_uid, m.content, m.created_at 
        FROM messages m
        WHERE m.chatroom_id = ?
        ORDER BY m.created_at ASC`,
		chatroomID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var createdAt string

		if err := rows.Scan(&msg.UserID, &msg.Message, &createdAt); err != nil {
			return nil, err
		}

		msg.CreatedAt = createdAt
		messages = append(messages, msg)
	}

	return messages, nil
}
