package _interface

import (
	"context"
)

type IGetUserRepository interface {
	FindOneUser(ctx context.Context, uID uint) error
}
