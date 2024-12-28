package entity

import "mime/multipart"

type UpdateProfileUserEntity struct {
	Image  *multipart.FileHeader `json:"image"`
	UserID uint                  `json:"userID"`
}
