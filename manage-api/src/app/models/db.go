package models

import (
	"fmt"
	"winddies-api/src/app/global"

	"github.com/go-redis/redis"
)

var RedisDb *redis.Client

func Init() {
	conf := global.Conf.Redis
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,x
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := RedisDb.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connect redis success...")
	}
}
