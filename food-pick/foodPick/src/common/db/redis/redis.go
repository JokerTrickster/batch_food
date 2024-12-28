package _redis

import (
	"context"
	"fmt"
	"log"
	"os"

	_aws "main/common/aws"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

const RankingKey = "food:rankings"
const PrevRankingKey = "prev:food:rankings"
const DailyKey = "daily:foods"
const NewRank = 0

func InitRedis() error {
	ctx := context.Background()
	isLocal := os.Getenv("IS_LOCAL")
	var connectionString string
	if isLocal == "true" {
		connectionString = fmt.Sprintf("redis://%s:%s@localhost:6379/0", os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"))
	} else {
		dbInfos, err := _aws.AwsSsmGetParams([]string{"dev_food_redis_user", "dev_food_redis_password", "dev_common_redis_host", "dev_common_redis_port", "dev_food_redis_db"})
		if err != nil {
			return err
		}
		fmt.Println(dbInfos)
		connectionString = fmt.Sprintf("redis://%s:%s@%s:%s/%s",
			dbInfos[4], //user
			dbInfos[3], //password
			dbInfos[0], //host
			dbInfos[1], //port
			dbInfos[2], //db

		)
		fmt.Println(connectionString)
	}

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		log.Println(err)
		return err
	}

	Client = redis.NewClient(opt)

	_, err = Client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("Connected to Redis!")

	return nil
}
