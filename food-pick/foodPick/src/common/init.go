package utils

import (
	"fmt"
	"main/common/aws"
	_aws "main/common/aws"
	"main/common/db/mysql"
	_redis "main/common/db/redis"
	"os"
)

var GeminiID string

func InitServer() error {
	if err := InitEnv(); err != nil {
		fmt.Sprintf("서버 에러 발생 : %s", err.Error())
		return err
	}

	if err := InitJwt(); err != nil {
		fmt.Sprintf("jwt 초기화 에러 : %s", err.Error())
		return err
	}
	if err := _aws.InitAws(); err != nil {
		fmt.Sprintf("aws 초기화 에러 : %s", err.Error())
		return err
	}

	if err := InitGoogleOauth(); err != nil {
		fmt.Sprintf("google oauth 초기화 에러 : %s", err.Error())
		return err
	}
	if err := InitKakaoOauth(); err != nil {
		fmt.Sprintf("kakao oauth 초기화 에러 : %s", err.Error())
		return err
	}
	if err := mysql.InitMySQL(); err != nil {
		fmt.Sprintf("db 초기화 에러 : %s", err.Error())
		return err
	}
	if err := _redis.InitRedis(); err != nil {
		fmt.Sprintf("redis 초기화 에러 : %s", err.Error())
		return err
	}
	if err := InitNotice(); err != nil {
		fmt.Sprintf("notice 초기화 에러 : %s", err.Error())
		return err
	}

	if Env.IsLocal {
		var ok bool
		GeminiID, ok = os.LookupEnv("GEMINI_ID")
		if !ok {
			fmt.Println("GEMINI_ID not found")
		}

	} else {
		var err error
		GeminiID, err = aws.AwsSsmGetParam("food_gemini_id")
		if err != nil {
			fmt.Println(err)
		}
		if err := InitLogging(); err != nil {
			return err
		}
	}

	return nil
}
