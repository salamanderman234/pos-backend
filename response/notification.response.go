package response

type NotificationResponse struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Date    int64  `json:"date"`
	IsRead  bool   `json:"is_read"`
}
