package common

import (
	"common/aws"
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"google.golang.org/api/option"
)

var MessageClient *messaging.Client

func InitNotice() error {
	ctx := context.Background()
	serviceKey, err := aws.AwsSsmGetParam("firebase_service_key")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 서비스 계정 JSON 키를 byte 배열로 변환합니다.
	credentials := []byte(serviceKey)
	opt := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return ErrorMsg(ctx, ErrInternalServer, Trace(), fmt.Sprintf("failed to initialize app: %v", err), ErrFromFirebase)
	}
	MessageClient, err = app.Messaging(ctx)
	if err != nil {
		return ErrorMsg(ctx, ErrInternalServer, Trace(), fmt.Sprintf("error getting Messaging client: %v", err), ErrFromFirebase)
	}

	return nil
}
