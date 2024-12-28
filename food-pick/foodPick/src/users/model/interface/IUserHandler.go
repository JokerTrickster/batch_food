package _interface

import (
	"github.com/aws/aws-lambda-go/events"
)

type IGetUserHandler interface {
	GetUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
