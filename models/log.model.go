package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	LogType string `json:"log_type"`
	Data    string `json:"data"`
	Message string `json:"message"`
	Date    int64  `json:"date"`
}
