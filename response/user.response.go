package response

type UserResponse struct {
	ID                 uint                   `json:"id"`
	Level              string                 `json:"level"`
	Status             string                 `json:"status"`
	Username           string                 `json:"username"`
	Email              string                 `json:"email"`
	Fullname           string                 `json:"fullname"`
	Province           string                 `json:"province"`
	City               string                 `json:"city"`
	FullAddress        string                 `json:"full_address"`
	IsTwoFactorEnabled bool                   `json:"is_two_factor_enabled"`
	TwoFactorMethod    string                 `json:"two_factor_method"`
	VerifiedAt         int64                  `json:"verified_at"`
	LastChangePassword int64                  `json:"last_change_password"`
	Notifications      []NotificationResponse `json:"notifications"`
	Devices            []UserDeviceResponse   `json:"devices"`
}

type UserDeviceResponse struct {
	ID           uint   `json:"id"`
	Device       string `json:"device"`
	Type         string `json:"type"`
	LastLogin    int64  `json:"last_login"`
	LastActivity int64  `json:"last_activity"`
	BannedAt     int64  `json:"banned_at"`
	BanReason    string `json:"ban_reason"`
	BannedBy     string `json:"banned_by"`
}
