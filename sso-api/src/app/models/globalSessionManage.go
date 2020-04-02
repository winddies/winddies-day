package models

import (
	"encoding/json"
)

const clientHashKey = "sso_client"

func NewGlobalSessionManage() GlobalSession {
	return &globalSessionManage{}
}

type GlobalSessionManageInfo struct {
	ClientId   string `mapstructure:"clientId"`
	LogoutUrl  string `mapstructure:"logoutUrl"`
	ClientName string `mapstructure:"clientName"`
}

type GlobalSession interface {
	Set(globalSessionId string, info *GlobalSessionManageInfo) error
	GetClientIds(globalSessionId string) ([]string, error)
	GetClientInfo(clientId string) (result map[string]interface{}, err error)
	DeleteGlobalSession(globalSessionId string) error
}

type globalSessionManage struct{}

func (globalSession *globalSessionManage) Set(globalSessionId string, info *GlobalSessionManageInfo) error {

	_, err := RedisDb.LPush(globalSessionId, info.ClientId).Result()
	if err != nil {
		return err
	}

	var data map[string]interface{}
	// t := reflect.TypeOf(info)
	// v := reflect.ValueOf(info)
	// for i := 0; i < t.NumField(); i++ {
	// 	data[t.Field(i).Name] = v.Field(i).Interface()
	// }
	b, _ := json.Marshal(info)
	json.Unmarshal(b, &data)

	_, err = RedisDb.HSet(clientHashKey, info.ClientId, data).Result()

	return err
}

func (globalSession *globalSessionManage) GetClientIds(globalSessionId string) ([]string, error) {
	result, err := RedisDb.LRange(globalSessionId, 0, -1).Result()
	return result, err
}

func (GlobalSession *globalSessionManage) DeleteGlobalSession(globalSessionId string) error {
	_, err := RedisDb.Del("session_" + globalSessionId).Result()
	return err
}

func (globalSession *globalSessionManage) GetClientInfo(clientId string) (result map[string]interface{}, err error) {
	value, err := RedisDb.HGet(clientHashKey, clientId).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, err
	}

	return
}
