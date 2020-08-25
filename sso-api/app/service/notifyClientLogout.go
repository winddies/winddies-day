package service

import (
	"io/ioutil"
	"luck-home/winddies/sso-api/app/models"
	"net/http"
	"net/url"
	"sync"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func NotifyClientLogout(logger *log.Entry, clientIds []string, GlobalSessionManage models.GlobalSession) {
	var wg sync.WaitGroup
	for _, val := range clientIds {
		wg.Add(1)
		go func(key string) {
			clientLogout(logger, GlobalSessionManage, key)
			wg.Done()
		}(val)
	}
	wg.Wait()
}

func clientLogout(logger *log.Entry, GlobalSessionManage models.GlobalSession, key string) {
	mapData, err := GlobalSessionManage.GetClientInfo(key)
	if err != nil {
		logger.Errorf("<GetClientInfo> GlobalSessionManage get clientInfo error, %s", err)
		return
	}

	var info *models.GlobalSessionManageInfo

	if err := mapstructure.Decode(mapData, info); err != nil {
		logger.Errorf("<mapstructure.Decode> decode error, %s", err)
		return
	}

	base, err := url.Parse(info.LogoutUrl)
	if err != nil {
		return
	}
	// Query params
	params := url.Values{}
	params.Add("sessionId", info.ClientId)
	base.RawQuery = params.Encode()
	client := &http.Client{}

	req, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		logger.Errorf("<http.NewRequest> client %s logout error, %s", info.ClientName, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("<ioutil.ReadAll>clientId %s read response body ereror, %s", info.ClientId, err)
	}

	defer resp.Body.Close()

	logger.Info(string(body))
}
