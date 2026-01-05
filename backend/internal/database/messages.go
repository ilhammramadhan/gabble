package database

import (
	"context"

	"github.com/ilhammramadhan/gabble/internal/models"
)

func (db *DB) CreateMessage(ctx context.Context, roomID, userID, content string) (*models.Message, error) {
	var msg models.Message
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO messages (room_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, room_id, user_id, content, created_at
	`, roomID, userID, content).Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (db *DB) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]models.Message, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT m.id, m.room_id, m.user_id, m.content, m.created_at,
			   u.id, u.github_id, u.username, u.avatar_url, u.created_at
		FROM messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at ASC
		LIMIT $2 OFFSET $3
	`, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var user models.User
		if err := rows.Scan(
			&msg.ID, &msg.RoomID, &msg.UserID, &msg.Content, &msg.CreatedAt,
			&user.ID, &user.GithubID, &user.Username, &user.AvatarURL, &user.CreatedAt,
		); err != nil {
			return nil, err
		}
		msg.User = &user
		messages = append(messages, msg)
	}
	return messages, nil
}
