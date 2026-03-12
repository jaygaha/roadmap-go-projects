package models

import "time"

type Score struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	GameID      string    `json:"game_id" gorm:"not null"`
	Score       int       `json:"score" gorm:"not null"`
	SubmittedAt time.Time `json:"submitted_at"`
}
