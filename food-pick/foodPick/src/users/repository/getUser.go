package repository

import (
	"context"
	"main/common/db/mysql"
	_interface "main/users/model/interface"

	"gorm.io/gorm"
)

func NewGetUserRepository(gormDB *gorm.DB) _interface.IGetUserRepository {
	return &GetUserRepository{GormDB: gormDB}
}
func (d *GetUserRepository) FindOneUser(ctx context.Context, uID uint) (*mysql.Users, error) {

	user := mysql.Users{}

	return &user, nil
}
