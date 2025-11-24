package models

import "time"

type Question struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `json:"created_at"`

	Answers []Answer `gorm:"constraint:OnDelete:CASCADE;" json:"answers,omitempty"`
}

type Answer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"not null;index" json:"question_id"`
	UserID     string    `gorm:"type:uuid;not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}
