package models

import "time"

type Session struct {
	SessionId  string    `json:"session_id"`
	InstanceId string 	  `json:"instance_id"`
	UserId     string    `json:"user_id"`
	GameName   string    `json:"game_name"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	ServerIP   string    `json:"server_ip"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	WorkSpace  string    `json:"workspace"`
	InstanceId string    `json:"instance_id"`
}
