package common

import (
	"fmt"
	"os"
	"time"
)

type envStruct struct {
	Port               string
	Env                string
	IsLocal            bool
	GoogleClientID     string
	GoogleClientSecret string
}

// Env : Environment
var Env envStruct

// 사용하는 환경 변수 네임 설정 함수
func InitVarNames() []string {
	result := make([]string, 0)
	result = append(result, "PORT")
	result = append(result, "ENV")
	result = append(result, "IS_LOCAL")
	return result
}

type CtxValues struct {
	Method    string
	Url       string
	UserID    uint
	StartTime time.Time
	RequestID string
	Email     string
}

// 사용할 환경 변수 값들 초기화해주는 함수
func InitEnv() error {
	envVarNames := InitVarNames()
	envs, err := getOSLookupEnv(envVarNames)
	if err != nil {
		return err
	}
	Env = envStruct{
		Port:    envs["PORT"],
		Env:     envs["ENV"],
		IsLocal: envIsLocal(envs["IS_LOCAL"]),
	}
	return nil
}

func envIsLocal(isLocal string) bool {
	if isLocal != "true" {
		return false
	} else {
		return true
	}
}
func getOSLookupEnv(envVarNames []string) (map[string]string, error) {
	result := map[string]string{}
	var ok bool
	for _, envVarName := range envVarNames {
		if result[envVarName], ok = os.LookupEnv(envVarName); !ok {
			return nil, fmt.Errorf("os lookup get failed")
		}
	}
	return result, nil
}

func TimeToEpochMillis(time time.Time) int64 {
	nanos := time.UnixNano()
	millis := nanos / 1000000
	return millis
}

func EpochToTime(date int64) time.Time {
	return time.Unix(date, 0)
}

func EpochToTimeMillis(t int64) time.Time {
	return time.Unix(t/1000, t%1000*1000000)
}

func PtrTrue() *bool {
	b := true
	return &b
}

func PtrFalse() *bool {
	b := false
	return &b
}
