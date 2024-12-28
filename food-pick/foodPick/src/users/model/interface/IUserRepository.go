package _interface

import (
	"context"
	"main/common/db/mysql"
)

type IGetUserRepository interface {
	FindOneUser(ctx context.Context, uID uint) (*mysql.Users, error)
}
