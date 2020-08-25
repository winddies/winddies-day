package models

import (
	"fmt"
	"luck-home/winddies/sso-api/app/global"

	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/go-redis/redis"
)

var RedisDb *redis.Client

var RedisSessionStore redisStore.Store

func Init() {
	conf := global.Conf.Redis
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := RedisDb.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connect redis success...")
	}

	RedisSessionStore, _ = redisStore.NewStore(10, "tcp", global.Conf.Redis.Addr, global.Conf.Redis.Password, []byte(global.Conf.Session.Secret))
}
