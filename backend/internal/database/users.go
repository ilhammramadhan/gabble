package database

import (
	"context"

	"github.com/ilhammramadhan/gabble/internal/models"
)

func (db *DB) CreateUser(ctx context.Context, githubID, username, avatarURL string) (*models.User, error) {
	var user models.User
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO users (github_id, username, avatar_url)
		VALUES ($1, $2, $3)
		ON CONFLICT (github_id) DO UPDATE SET
			username = EXCLUDED.username,
			avatar_url = EXCLUDED.avatar_url
		RETURNING id, github_id, username, avatar_url, created_at
	`, githubID, username, avatarURL).Scan(
		&user.ID, &user.GithubID, &user.Username, &user.AvatarURL, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := db.Pool.QueryRow(ctx, `
		SELECT id, github_id, username, avatar_url, created_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.GithubID, &user.Username, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByGithubID(ctx context.Context, githubID string) (*models.User, error) {
	var user models.User
	err := db.Pool.QueryRow(ctx, `
		SELECT id, github_id, username, avatar_url, created_at
		FROM users WHERE github_id = $1
	`, githubID).Scan(&user.ID, &user.GithubID, &user.Username, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
