package main

import (
	"fmt"
	"main/common/db/mysql"
	"main/users/handler"
	"main/users/repository"
	"main/users/usecase"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(req.Path)
	fmt.Println("테스트 중입니다.d")
	switch req.Path {
	case "/v0.1/users":
		return handler.NewGetUserHandler(usecase.NewGetUserUseCase(repository.NewGetUserRepository(mysql.GormMysqlDB), mysql.DBTimeOut)).GetUser(req)
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
