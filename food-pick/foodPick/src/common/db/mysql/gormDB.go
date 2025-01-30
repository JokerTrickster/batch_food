package mysql

import "gorm.io/gorm"

// 시나리오 상수 정의
// 연인, 혼반, 가족, 다이어트, 회식, 친구
const (
	ScenarioAll    = 0
	ScenarioCouple = iota
	ScenarioSolo
	ScenarioFamily
	ScenarioCompany
	ScenarioFriend
)

// 식사 시간 상수 정의
// 아침, 점심, 저녁, 브런치, 간식, 야식
const (
	TimeAll      = 0
	TimeMorning  = iota //아침
	TimeLunch           //점심
	TimeDinner          //저녁
	TimeSnack           //간식
	TimeMidnight        //야식
)

// 음식 종류 상수 정의
// 한식, 중식, 일식, 양식, 분식,베트남 음식, 인도 음식, 패스트 푸드, 디저트, 퓨전 요리
const (
	TypeAll        = 0
	TypeKorean     = iota //한식
	TypeChinese           //중식
	TypeJapanese          //일식
	TypeWestern           //양식
	TypeStreetFood        //분식
	TypeFastFood          //패스트 푸드
	TypeVietnamese        //베트남 음식
	TypeIndian            //인도 음식
	TypeDessert           //디저트
	TypeFusion            //퓨전 요리
)

// 기분/테마 상수 정의
// 스트레스 해소, 피로 회복, 기분 전환, 제철 음식, 영양식, 특별한 날
const (
	ThemeAll             = 0
	ThemeStressRelief    = iota // 스트레스 해소
	ThemeHangover               // 해장
	ThemeFatigueRecovery        // 피로 회복
	ThemeDiet                   //다이어트
	ThemeSeasonalFood           //제철 음식
)

// 맵 정의
var ScenarioMap = map[string]int{
	"연인": ScenarioCouple,
	"혼밥": ScenarioSolo,
	"가족": ScenarioFamily,
	"회식": ScenarioCompany,
	"친구": ScenarioFriend,
}

var TimeMap = map[string]int{
	"아침": TimeMorning,
	"점심": TimeLunch,
	"저녁": TimeDinner,
	"간식": TimeSnack,
	"야식": TimeMidnight,
}

var TypeMap = map[string]int{
	"한식":     TypeKorean,
	"중식":     TypeChinese,
	"일식":     TypeJapanese,
	"양식":     TypeWestern,
	"분식":     TypeStreetFood,
	"베트남 음식": TypeVietnamese,
	"인도 음식":  TypeIndian,
	"패스트 푸드": TypeFastFood,
	"디저트":    TypeDessert,
	"퓨전 요리":  TypeFusion,
}

var ThemeMap = map[string]int{
	"스트레스 해소": ThemeStressRelief,
	"해장":      ThemeHangover,
	"피로 회복":   ThemeFatigueRecovery,
	"다이어트":    ThemeDiet,
	"제철 음식":   ThemeSeasonalFood,
}

// 반대 맵 정의 (int -> string)
var ScenarioReverseMap = make(map[int]string)
var TimeReverseMap = make(map[int]string)
var TypeReverseMap = make(map[int]string)
var ThemeReverseMap = make(map[int]string)

type Tokens struct {
	gorm.Model
	UserID           uint   `json:"userID" gorm:"column:user_id"`
	AccessToken      string `json:"accessToken" gorm:"column:access_token"`
	RefreshToken     string `json:"refreshToken" gorm:"column:refresh_token"`
	RefreshExpiredAt int64  `json:"refreshExpiredAt" gorm:"column:refresh_expired_at"`
}

type Users struct {
	gorm.Model
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
	Birth    string `json:"birth" gorm:"column:birth"`
	Name     string `json:"name" gorm:"column:name"`
	Sex      string `json:"sex" gorm:"column:sex"`
	Provider string `json:"provider" gorm:"column:provider"`
	Push     *bool  `json:"push" gorm:"column:push"`
	Image    string `json:"image" gorm:"column:image"`
}

type MetaTables struct {
	gorm.Model
	TableName        string `json:"tableName" gorm:"column:table_name"`
	TableDescription string `json:"tableDescription" gorm:"column:table_description"`
}

type Scenarios struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type Times struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type Types struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type Themes struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type UserAuths struct {
	gorm.Model
	Email    string `json:"email" gorm:"column:email"`
	AuthCode string `json:"authCode" gorm:"column:auth_code"`
	Type     string `json:"type" gorm:"column:type"`
}

type FoodImages struct {
	gorm.Model
	Name  string `json:"name" gorm:"column:name"`
	Image string `json:"image" gorm:"column:image"`
}

type Reports struct {
	gorm.Model
	UserID int    `json:"userID" gorm:"column:user_id"`
	Reason string `json:"reason" gorm:"column:reason"`
}

type Nutrients struct {
	gorm.Model
	FoodName     string  `json:"foodName" gorm:"column:food_name"`
	Amount       string  `json:"amount" gorm:"column:amount"`
	Kcal         float64 `json:"kcal" gorm:"column:kcal"`
	Carbohydrate float64 `json:"carbohydrate" gorm:"column:carbohydrate"`
	Protein      float64 `json:"protein" gorm:"column:protein"`
	Fat          float64 `json:"fat" gorm:"column:fat"`
}

type UserTokens struct {
	gorm.Model
	UserID uint   `json:"userID" gorm:"column:user_id"`
	Token  string `json:"token" gorm:"column:token"`
}

type Foods struct {
	gorm.Model
	Name    string `json:"name" gorm:"column:name"`
	ImageID int    `json:"imageID" gorm:"column:image_id"`
}
type CategoryTypes struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name"`
}

type Categories struct {
	gorm.Model
	Name   string `json:"name" gorm:"column:name"`
	TypeID int    `json:"typeID" gorm:"column:type_id"`
}

type FoodCategories struct {
	gorm.Model
	FoodID     int `json:"foodID" gorm:"column:food_id"`
	CategoryID int `json:"categoryID" gorm:"column:category_id"`
}

type FoodHistories struct {
	gorm.Model
	UserID int `json:"userID" gorm:"column:user_id"`
	FoodID int `json:"foodID" gorm:"column:food_id"`
}
