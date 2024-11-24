package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var globalErr error

func callAPI(ctx context.Context) error {
	// API URL
	apiURL := "https://dev-food-recommendation-api.jokertrickster.com/v0.1/system/foods/report"

	// 요청 데이터(JSON Body)
	requestBody := []byte(`{
		"key": "value"
	}`)

	// HTTP POST 요청 생성
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 요청 헤더 설정
	req.Header.Set("Content-Type", "application/json")

	// HTTP 클라이언트 생성
	client := &http.Client{}

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	// 응답 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// 응답 로그 출력
	log.Printf("API Response: %s\n", string(body))

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

func handler(ctx context.Context, request events.CloudWatchEvent) error {
	if globalErr != nil {
		return globalErr
	}

	// API 호출
	err := callAPI(ctx)
	if err != nil {
		log.Printf("Error calling API: %v", err)
		return err
	}

	log.Println("API call successful.")
	return nil
}

func main() {
	lambda.Start(handler)
}
