package models

import (
	"fmt"
	"login/src/app/utils"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/mitchellh/mapstructure"
)

type User struct {
	Password string `json:"password,omitempty"`
	UserName string `json:"userName"`
}

type UserInfo struct {
	UserName     string `mapstructure:"userName"`
	UserId       string `mapstructure:"userId"`
	Registertime string `mapstructure:"registertime"`
}

func (this *User) IsExist() (bool, error) {
	nameKey := "userName:" + strings.ToLower(strings.TrimSpace(this.UserName))
	cmd, err := RedisDb.Exists(nameKey).Result()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if 0 != cmd {
		return true, err
	}

	return false, err
}

func (this *User) IncrKey(incrName string) (incr_Key string) {
	result, err := RedisDb.Incr("incrCount").Result()
	if err != nil {
		panic(err)
	}

	incr_Key = incrName + ":" + strconv.FormatInt(result, 10)

	return incr_Key
}

func (this *User) Register() (bool, error) {
	exist, err := this.IsExist()
	if exist {
		return exist, err
	}

	uniqueKey := this.IncrKey("userId")
	stat, err := RedisDb.Set("userName:"+this.UserName, uniqueKey, 0).Result()
	if err != nil {
		return exist, err
	}

	if stat != "OK" {
		panic(stat)
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return exist, err
	}

	// Generate a snowflake ID.
	id := node.Generate()
	mess := map[string]interface{}{
		"userName":     this.UserName,
		"password":     utils.GetMd5(this.Password),
		"registertime": time.Now(),
		"userId":       int(id),
	}

	_, err = RedisDb.HMSet(uniqueKey, mess).Result()

	return exist, err
}

func (this *User) VerifyAuth() (bool, error) {
	exist, err := this.IsExist()
	if err != nil {
		return exist, err
	}

	if exist {
		uniqueKey, err := RedisDb.Get("userName:" + this.UserName).Result()
		if err != nil {
			return exist, err
		}

		result, err := RedisDb.HGet(uniqueKey, "password").Result()
		if err != nil {
			return exist, err
		}

		if result == utils.GetMd5(this.Password) {
			return exist, nil
		}
	}

	return exist, nil
}

func (this *User) GetUserInfo() (result *UserInfo, err error) {
	uniqueKey, err := RedisDb.Get("userName:" + this.UserName).Result()
	if err != nil {
		return
	}

	if uniqueKey == "" {
		return nil, nil
	}
	data, err := RedisDb.HGetAll(uniqueKey).Result()
	if err != nil {
		return nil, err
	}
	if err = mapstructure.Decode(data, &result); err != nil {
		return nil, err
	}
	return
}
