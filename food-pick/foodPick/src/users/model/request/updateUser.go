package request

type ReqUpdateUser struct {
	Birth        string `json:"birth" `
	Sex          string `json:"sex" `
	Name         string `json:"name" `
	NewPassword  string `json:"newPassword" `
	PrevPassword string `json:"prevPassword" `
	Push         *bool  `json:"push" `
}
