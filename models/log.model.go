package models

import (
	"gorm.io/gorm"
)

type LogLogin struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Device    string `json:"device"`
	Ip        string `json:"ip"`
	Date      int64  `json:"date"`
	IsSuccess bool   `json:"is_success"`
}
