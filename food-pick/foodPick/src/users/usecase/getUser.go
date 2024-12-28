package usecase

import (
	"context"
	"fmt"
	_interface "main/users/model/interface"
	"main/users/model/response"
	"time"
)

type GetUserUseCase struct {
	Repository     _interface.IGetUserRepository
	ContextTimeout time.Duration
}

func NewGetUserUseCase(repo _interface.IGetUserRepository, timeout time.Duration) _interface.IGetUserUseCase {
	return &GetUserUseCase{
		Repository:     repo,
		ContextTimeout: timeout,
	}
}

func (d *GetUserUseCase) Get(c context.Context, uID uint) (response.ResGetUser, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	fmt.Println(ctx)
	return response.ResGetUser{}, nil

}
