package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	// 1. 액세스 토큰 받아오기
	token, err := getAccessToken()
	if err != nil {
		fmt.Println("액세스 토큰 요청 실패:", err)
		return err
	}

	// 2. 음식 추천 API 호출
	if err := callRecommendAPI(token); err != nil {
		fmt.Println("음식 추천 API 요청 실패:", err)
		return err
	}

	return nil
}

func main() {

	lambda.Start(handler)
}

// 액세스 토큰을 가져오는 함수
func getAccessToken() (string, error) {
	url := "https://dev-food-recommendation-api.jokertrickster.com/v0.1/auth/guest"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("HTTP 요청 생성 실패: %v", err)
	}

	// HTTP 클라이언트 생성
	client := &http.Client{Timeout: 10 * time.Second}

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("응답 본문 읽기 실패: %v", err)
	}

	// 응답을 구조체로 변환
	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return "", fmt.Errorf("응답 파싱 실패: %v", err)
	}

	// 액세스 토큰 반환
	return authResponse.AccessToken, nil
}

func callRecommendAPI(token string) error {
	// 요청 데이터 생성
	data := generateRandomRecommendRequest()
	fmt.Println(data)
	// 요청 데이터를 JSON 형식으로 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON 변환 실패: %v", err)
	}

	// 음식 추천 API 호출
	url := "https://dev-food-recommendation-api.jokertrickster.com/v0.1/foods/recommend"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP 요청 생성 실패: %v", err)
	}

	// 요청 헤더 설정 (tkn 헤더에 액세스 토큰 추가)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("tkn", token)

	// HTTP 클라이언트 생성
	client := &http.Client{Timeout: 10 * time.Second}

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP 요청 실패: 상태 코드 %d", resp.StatusCode)
	}
	// 응답 본문 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("응답 본문 읽기 실패: %v", err)
	}

	// 응답 데이터 구조체로 변환
	var recommendResponse RecommendResponse
	if err := json.Unmarshal(body, &recommendResponse); err != nil {
		return fmt.Errorf("응답 데이터 변환 실패: %v", err)
	}

	// 응답 출력 (음식 이름과 이미지)
	for _, food := range recommendResponse.FoodNames {
		fmt.Printf("음식 이름: %s, 이미지 URL: %s\n", food.Name)
	}
	fmt.Println("음식 추천 API 호출 성공")
	return nil
}
