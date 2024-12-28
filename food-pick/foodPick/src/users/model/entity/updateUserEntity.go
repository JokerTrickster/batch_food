package entity

type UpdateUserEntity struct {
	UserID       uint   `json:"user_id"`
	Birth        string `json:"birth"`
	Name         string `json:"name"`
	Sex          string `json:"sex"`
	Email        string `json:"email"`
	PrevPassword string `json:"prevPassword"`
	NewPassword  string `json:"newPassword"`
	Push         *bool  `json:"push"`
}
