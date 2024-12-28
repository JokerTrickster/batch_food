package response

type ResGetUser struct {
	Name   string `json:"name"`
	Birth  string `json:"birth"`
	Sex    string `json:"sex"`
	Email  string `json:"email"`
	Push   *bool  `json:"push"`
	Image  string `json:"image"`
	UserID int    `json:"userID"`
}
