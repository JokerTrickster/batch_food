package _interface

import (
	"context"
	"main/model/response"
)

type IGetUserUseCase interface {
	Get(c context.Context, uID uint) (response.ResGetUser, error)
}
