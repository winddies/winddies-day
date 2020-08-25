package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type AuthTokenInfo struct {
	UserId          string `mapstructure:"userId" json:"userId"`
	UserName        string `mapstructure:"userName" json:"userName"`
	ClientId        string `mapstructure:"clientId" json:"clientId"`
	GlobalSessionId string `mapstructure:"globalSessionId" json:"globalSessionId"`
}

type AuthToken interface {
	SetTokenInfo(data AuthTokenInfo) error
	GetTokenInfo() (map[string]string, error)
}

type authToken struct {
	tokenId string
}

func NewAuthToken(token string) AuthToken {
	return &authToken{tokenId: token}
}

func IncrKey(incrName string) (incr_Key string) {
	result, err := RedisDb.Incr("incrTokenId").Result()
	if err != nil {
		panic(err)
	}

	incr_Key = incrName + ":" + strconv.FormatInt(result, 10)

	return incr_Key
}

func getTokenKey(tokenId string) string {
	tokenKey := IncrKey("tokenId")
	RedisDb.Set(tokenId, tokenKey, 60*time.Second)
	return tokenKey
}

func (token *authToken) DeleteTokenInfo() (cmd int64, err error) {
	tokenKey, err := RedisDb.Get(token.tokenId).Result()
	if err != nil {
		return 0, err
	}

	cmd, err = RedisDb.HDel(tokenKey).Result()
	return cmd, err
}

func (token *authToken) SetTokenInfo(data AuthTokenInfo) error {
	tokenKey := getTokenKey(token.tokenId)
	var dataMap map[string]interface{}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &dataMap)

	_, err := RedisDb.HMSet(tokenKey, dataMap).Result()
	result, errs := RedisDb.HGetAll(tokenKey).Result()
	if errs == nil {
		fmt.Printf("%+v", result)
	}
	RedisDb.Expire(tokenKey, 60*time.Second)
	return err
}

func (token *authToken) GetTokenInfo() (result map[string]string, err error) {
	tokenKey, err := RedisDb.Get(token.tokenId).Result()
	fmt.Println(tokenKey)
	if err != nil {
		return nil, err
	}

	result, err = RedisDb.HGetAll(tokenKey).Result()
	return result, err
}
