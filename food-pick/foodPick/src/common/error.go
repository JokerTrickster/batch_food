package utils

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// 프론트엔드 받을 에러 형식
type ResError struct {
	ErrType string `json:"errType,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

// 에러 로깅을 위한 에러 형식
type Err struct {
	HttpCode int    `json:"httpCode,omitempty"`
	ErrType  string `json:"errType,omitempty"`
	Msg      string `json:"msg,omitempty"`
	Trace    string `json:"trace,omitempty"`
	From     string `json:"from,omitempty"`
}

// 에러 타입을 구분
type ErrType string

// 에러가 어디서 발생했는지 확인용
type IErrFrom string

const (
	ErrFromClient   = IErrFrom("client")
	ErrFromInternal = IErrFrom("internal")
	ErrFromMongoDB  = IErrFrom("mongoDB")
	ErrFromMysqlDB  = IErrFrom("mysqlDB")
	ErrFromRedis    = IErrFrom("redis")
	ErrFromAws      = IErrFrom("aws")
	ErrFromAwsS3    = IErrFrom("aws_s3")
	ErrFromAwsSsm   = IErrFrom("aws_ssm")
	ErrFromNaver    = IErrFrom("naver")
	ErrFromGemini   = IErrFrom("gemini")
	ErrFromKakao    = IErrFrom("kakao")
	ErrFromFirebase = IErrFrom("firebase")
	ErrFromChatGPT  = IErrFrom("chatGPT")
)

// basic error
const (
	ErrBadParameter   = ErrType("PARAM_BAD")
	ErrNotFound       = ErrType("NOT_FOUND")
	ErrBadToken       = ErrType("TOKEN_BAD")
	ErrInternalServer = ErrType("INTERNAL_SERVER")
	ErrInternalDB     = ErrType("INTERNAL_DB")
	ErrPartner        = ErrType("PARTNER")
)

// auth error
const (
	ErrCodeNotFound           = ErrType("CODE_NOT_FOUND")
	ErrUserNotFound           = ErrType("USER_NOT_FOUND")
	ErrProfileNotFount        = ErrType("PROFILE_NOT_FOUND")
	ErrUserAlreadyExisted     = ErrType("USER_ALREADY_EXISTED")
	ErrInvalidAccessToken     = ErrType("INVALID_ACCESS_TOKEN")
	ErrPasswordNotMatch       = ErrType("PASSWORD_NOT_MATCH")
	ErrInvalidAuthCode        = ErrType("INVALID_AUTH_CODE")
	ErrInvalidEmailOrPassword = ErrType("INVALID_EMAIL_OR_PASSWORD") //패스워드 또는 이메일이 잘못됐습니다.
)

// food error
const (
	ErrGeminiError  = ErrType("GEMINI_INTERNAL_SERVER")
	ErrFoodNotFound = ErrType("FOOD_NOT_FOUND")
)

// basic , game, room, auth error mapping
var ErrHttpCode = map[string]int{
	//400
	"PARAM_BAD":                 http.StatusBadRequest,
	"USER_ALREADY_EXISTED":      http.StatusBadRequest,
	"BAD_REQUEST":               http.StatusBadRequest,
	"USER_NOT_FOUND":            http.StatusBadRequest,
	"NOT_ENOUGH_CARD":           http.StatusBadRequest,
	"NOT_ENOUGH_CONDITION":      http.StatusBadRequest,
	"PASSWORD_NOT_MATCH":        http.StatusBadRequest,
	"INVALID_EMAIL_OR_PASSWORD": http.StatusBadRequest,
	"FOOD_NOT_FOUND":            http.StatusBadRequest,

	//401 인증이 필요한 경우 (인증)
	"TOKEN_BAD":            http.StatusUnauthorized,
	"INVALID_ACCESS_TOKEN": http.StatusUnauthorized,
	"INVALID_AUTH_CODE":    http.StatusUnauthorized,

	//403 인증은 됐으나 권한이 없는 경우 (인가)
	"PARTNER": http.StatusForbidden,

	//404
	"NOT_FOUND": http.StatusNotFound,

	//500
	"INTERNAL_SERVER":        http.StatusInternalServerError,
	"INTERNAL_DB":            http.StatusInternalServerError,
	"GEMINI_INTERNAL_SERVER": http.StatusInternalServerError,
}

func ErrorParsing(data string) Err {
	slice := strings.Split(data, "|")
	result := Err{
		HttpCode: ErrHttpCode[slice[0]],
		ErrType:  slice[0],
		Trace:    slice[1],
		Msg:      slice[2],
		From:     slice[3],
	}
	return result
}

func ErrorMsg(ctx context.Context, errType ErrType, trace string, msg string, from IErrFrom) error {

	return fmt.Errorf("%s|%s|%s|%s", errType, trace, msg, from)
}

func (e ErrType) New(errType string, msg string) *ResError {
	return &ResError{ErrType: errType, Msg: msg}
}

func Trace() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	_, line := runtime.FuncForPC(pc).FileLine(pc)
	return fmt.Sprintf("%s.L%d", funcName, line)
}

func HandleError(errMsg string, args ...interface{}) string {
	// 인자로 받은 값들을 문자열로 변환
	var params []string
	for _, arg := range args {
		params = append(params, fmt.Sprintf("%v", arg))
	}

	// 에러 메시지와 파라미터들을 조합
	result := fmt.Sprintf("Error: %s | Parameters: %s", errMsg, strings.Join(params, ", "))

	return result
}
