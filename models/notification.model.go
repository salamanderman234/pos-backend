package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Date    int64  `json:"date"`
	IsRead  bool   `json:"is_read"`
}
