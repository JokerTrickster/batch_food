package mysql

import (
	"database/sql"
	"fmt"
	_aws "main/common/aws"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *sql.DB
var GormMysqlDB *gorm.DB

const DBTimeOut = 8 * time.Second

func InitMySQL() error {
	var connectionString string
	var err error
	isLocal := os.Getenv("IS_LOCAL")
	if isLocal == "true" {
		// MySQL 연결 문자열
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		)
	} else {
		dbInfos, err := _aws.AwsSsmGetParams([]string{"dev_food_mysql_user", "dev_food_mysql_password", "dev_common_mysql_host", "dev_common_mysql_port", "dev_food_mysql_db"})
		if err != nil {
			return err
		}
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbInfos[4], //user
			dbInfos[3], //host
			dbInfos[0], //db name,
			dbInfos[1], //port
			dbInfos[2], //db name
		)
	}
	fmt.Println(connectionString)
	// MySQL에 연결
	MysqlDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Failed to connect to MySQL!")
		fmt.Sprintln("에러 메시지 %s", err)
	}
	fmt.Println("Connected to MySQL!")

	/*
		GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		you can disable it during initialization if it is not required, you will gain about 30%+ performance improvement after that
	*/
	GormMysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: MysqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("Failed to connect to Gorm MySQL!")
		fmt.Sprintln("에러 메시지 %s", err)
	}

	// gen 패키지를 사용하여 쿼리를 생성할 때 사용할 DB를 설정
	// SetDefault(GormMysqlDB)
	InitMeta()

	return nil
}

func PKIDGenerate() string {
	//uuid 로 생성
	result := (uuid.New()).String()
	return result
}

func NowDateGenerate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func EpochToTime(t int64) time.Time {
	return time.Unix(t, t%1000*1000000)
}
func EpochToTimeString(t int64) string {
	return time.Unix(t, t%1000*1000000).String()
}

func TimeStringToEpoch(t string) int64 {
	date, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", t)
	return date.Unix()
}

func TimeToEpoch(t time.Time) int64 {
	return t.Unix()
}

func GetScenarioID(name string) (int, error) {
	id, exists := ScenarioMap[name]
	if !exists {
		return 0, fmt.Errorf("시나리오 이름을 찾을 수 없습니다: %s", name)
	}
	return id, nil
}

func GetThemeID(name string) (int, error) {
	id, exists := ThemeMap[name]
	if !exists {
		return 0, fmt.Errorf("테마 이름을 찾을 수 없습니다: %s", name)
	}
	return id, nil
}

func GetTimeID(name string) (int, error) {
	id, exists := TimeMap[name]
	if !exists {
		return 0, fmt.Errorf("시간 이름을 찾을 수 없습니다: %s", name)
	}
	return id, nil
}

func GetTypeID(name string) (int, error) {
	id, exists := TypeMap[name]
	if !exists {
		return 0, fmt.Errorf("음식 종류 이름을 찾을 수 없습니다: %s", name)
	}
	return id, nil
}

// 맵에서 키에 해당하는 값을 가져오는 함수
func GetScenarioKey(val int) (string, bool) {
	key, ok := ScenarioReverseMap[val]
	return key, ok
}

func GetTimeKey(val int) (string, bool) {
	key, ok := TimeReverseMap[val]
	return key, ok
}

func GetTypeKey(val int) (string, bool) {
	key, ok := TypeReverseMap[val]
	return key, ok
}

func GetThemeKey(val int) (string, bool) {
	key, ok := ThemeReverseMap[val]
	return key, ok
}

func InitMeta() {
	for k, v := range ScenarioMap {
		ScenarioReverseMap[v] = k
	}
	for k, v := range TimeMap {
		TimeReverseMap[v] = k
	}
	for k, v := range TypeMap {
		TypeReverseMap[v] = k
	}
	for k, v := range ThemeMap {
		ThemeReverseMap[v] = k
	}

}

// 트랜잭션 처리 미들웨어
func Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	if err = tx.Error; err != nil {
		return err
	}

	err = fc(tx)
	return
}
