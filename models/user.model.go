package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Avatar             string             `json:"avatar"`
	Level              string             `json:"level"`
	Status             string             `json:"status"`
	Username           string             `json:"username"`
	Email              string             `json:"email"`
	Password           string             `json:"password"`
	Fullname           string             `json:"fullname"`
	Province           string             `json:"province"`
	City               string             `json:"city"`
	FullAddress        string             `json:"full_address"`
	IsTwoFactorEnabled bool               `json:"is_two_factor_enabled"`
	TwoFactorMethod    string             `json:"two_factor_method"`
	Secret             string             `json:"secret"`
	Key                string             `json:"key"`
	KeyValidUntil      int64              `json:"key_valid_until"`
	KeyPurpose         string             `json:"key_purpose"`
	VerifiedAt         int64              `json:"verified_at"`
	BannedAt           int64              `json:"is_banned"`
	BannedBy           string             `json:"banned_by"`
	BanReason          string             `json:"ban_reason"`
	SuspendedAt        int64              `json:"suspended_at"`
	SuspendedUntil     int64              `json:"suspended_until"`
	SuspendedBy        string             `json:"suspended_by"`
	SuspendReason      string             `json:"suspend_reason"`
	LastChangePassword int64              `json:"last_change_password"`
	Notifications      []Notification     `json:"notifications"`
	Passwords          []UserPasswordHash `json:"passwords"`
}

type UserPasswordHash struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Hash   string `json:"hash"`
	Date   int64  `json:"date"`
}
