package main

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type RecommendRequest struct {
	Flavors        string `json:"flavors"`
	PreviousAnswer string `json:"previousAnswer"`
	Scenarios      string `json:"scenarios"`
	Themes         string `json:"themes"`
	Times          string `json:"times"`
	Types          string `json:"types"`
}

// 응답 데이터 구조체 정의
type RecommendResponse struct {
	FoodNames []Food `json:"foodNames"`
}

type Food struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

// 시나리오 상수 정의
var scenarios = []string{"연인", "혼밥", "가족", "다이어트", "회식", "친구"}

// 식사 시간 상수 정의
var times = []string{"아침", "점심", "저녁", "브런치", "간식", "야식"}

// 음식 종류 상수 정의
var types = []string{"한식", "중식", "일식", "양식", "분식", "베트남 음식", "인도 음식", "패스트 푸드", "디저트", "퓨전 요리"}

// 기분/테마 상수 정의
var themes = []string{"스트레스 해소", "피로 회복", "기분 전환", "제철 음식", "영양식", "특별한 날"}

// 맛 상수 정의
var flavors = []string{"매운맛", "감칠맛", "고소한맛", "단맛", "짠맛", "싱거운맛"}
