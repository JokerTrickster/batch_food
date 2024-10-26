package main

import "golang.org/x/exp/rand"

// 시나리오 상수 정의
// 연인, 혼반, 가족, 다이어트, 회식, 친구

// 식사 시간 상수 정의
// 아침, 점심, 저녁, 브런치, 간식, 야식

// 음식 종류 상수 정의
// 한식, 중식, 일식, 양식, 분식,베트남 음식, 인도 음식, 패스트 푸드, 디저트, 퓨전 요리

// 기분/테마 상수 정의
// 스트레스 해소, 피로 회복, 기분 전환, 제철 음식, 영양식, 특별한 날

// 맛 상수 정의
// 매운맛, 감칠맛, 고소한맛, 단맛, 짠맛, 싱거운맛

// 랜덤 값 생성 함수
func getRandomValue(options []string) string {
	// 빈 값일 확률 추가 (약 20%)
	if rand.Float32() < 0.2 {
		return ""
	}

	// 랜덤으로 배열에서 선택
	index := rand.Intn(len(options))
	return options[index]
}

// RecommendRequest 생성 함수
func generateRandomRecommendRequest() RecommendRequest {
	// 랜덤으로 필드 값 생성
	return RecommendRequest{
		Flavors:   getRandomValue(flavors),
		Scenarios: getRandomValue(scenarios),
		Themes:    getRandomValue(themes),
		Times:     getRandomValue(times),
		Types:     getRandomValue(types),
	}
}
