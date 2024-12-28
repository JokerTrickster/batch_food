package _interface

import (
	"context"
	"main/users/model/response"
)

type IGetUserUseCase interface {
	Get(c context.Context, uID uint) (response.ResGetUser, error)
}
