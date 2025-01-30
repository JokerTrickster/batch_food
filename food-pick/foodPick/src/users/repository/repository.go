package repository

import "gorm.io/gorm"

type GetUserRepository struct {
	GormDB *gorm.DB
}
