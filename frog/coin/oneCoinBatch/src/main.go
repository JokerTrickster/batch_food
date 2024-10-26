package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var globalErr error

func handler(ctx context.Context, request events.CloudWatchEvent) error {
	if globalErr != nil {
		return globalErr
	}
	// 1. api 호출
	err := callFullCoinAPI()
	if err != nil {
		fmt.Println("액세스 토큰 요청 실패:", err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

// 액세스 토큰을 가져오는 함수
func callFullCoinAPI() error {
	url := "https://dev-frog-api.jokertrickster.com/v2.1/users/batch/coins/one"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("HTTP 요청 생성 실패: %v", err)
	}

	// HTTP 클라이언트 생성
	client := &http.Client{Timeout: 10 * time.Second}

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("응답 본문 읽기 실패: %v", err)
	}

	return nil
}
