package models

import "time"

type Session struct {
	SessionId  string    `json:"session_id"`
	UserId     string    `json:"user_id"`
	GameName   string    `json:"game_name"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	InstanceID string    `json:"instance_id"`
	ServerIP   string    `json:"server_ip"`
	StateFile  string    `json:"state_file"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	WorkSpace  string    `json:"workspace"`
}
