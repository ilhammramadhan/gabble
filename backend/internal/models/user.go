package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	GithubID  string    `json:"github_id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}
