package handler

import (
	"fmt"
	_interface "main/users/model/interface"

	"github.com/aws/aws-lambda-go/events"
)

type GetUserHandler struct {
	UseCase _interface.IGetUserUseCase
}

func NewGetUserHandler(useCase _interface.IGetUserUseCase) _interface.IGetUserHandler {
	handler := &GetUserHandler{
		UseCase: useCase,
	}
	return handler
}

func (d *GetUserHandler) GetUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request
	fmt.Println("test")

	// Return response
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       "User created",
	}, nil
}
