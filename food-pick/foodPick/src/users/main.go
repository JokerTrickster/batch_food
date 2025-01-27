package main

import (
	"fmt"
	"main/handler"
	"main/repository"
	"main/usecase"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(req.Path)
	fmt.Println("테스트 중입니다.d")
	switch req.Path {
	case "/v0.1/users":
		return handler.NewGetUserHandler(usecase.NewGetUserUseCase(repository.NewGetUserRepository(nil), 8*time.Second)).GetUser(req)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("Path %s not found", req.Path),
		}, nil
	}
}

func main() {
	lambda.Start(router)
}
