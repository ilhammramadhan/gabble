package database

import (
	"context"

	"github.com/ilhammramadhan/gabble/internal/models"
)

func (db *DB) CreateRoom(ctx context.Context, name, createdBy string) (*models.Room, error) {
	var room models.Room
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO rooms (name, created_by)
		VALUES ($1, $2)
		RETURNING id, name, created_by, created_at
	`, name, createdBy).Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (db *DB) GetRooms(ctx context.Context) ([]models.Room, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, created_by, created_at
		FROM rooms
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (db *DB) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	var room models.Room
	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, created_by, created_at
		FROM rooms WHERE id = $1
	`, id).Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (db *DB) DeleteRoom(ctx context.Context, id, userID string) error {
	_, err := db.Pool.Exec(ctx, `
		DELETE FROM rooms WHERE id = $1 AND created_by = $2
	`, id, userID)
	return err
}
