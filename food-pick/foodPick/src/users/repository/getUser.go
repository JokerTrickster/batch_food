package repository

import (
	"context"
	_interface "main/model/interface"

	"gorm.io/gorm"
)

func NewGetUserRepository(gormDB *gorm.DB) _interface.IGetUserRepository {
	return &GetUserRepository{GormDB: gormDB}
}
func (d *GetUserRepository) FindOneUser(ctx context.Context, uID uint) error{

	return nil
}
