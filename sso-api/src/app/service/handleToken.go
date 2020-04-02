package service

import (
	"login/src/app/global"
	"login/src/app/models"
	"login/src/app/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 该函数用于处理 token 产生及 tokenInfo 的生成及保存
func HandleToken(ctx *gin.Context, logger *log.Entry, user *models.User) (token string, err error) {
	token = utils.GetUniqueToken()
	authToken := models.NewAuthToken(token)
	userInfo, err := user.GetUserInfo()
	logger.Infof("%+v", userInfo)
	if err != nil {
		logger.Errorf("<GetUserInfo> GetUserInfo error, %s", err)
		return
	}

	sessionStore, err := models.RedisSessionStore.Get(ctx.Request, global.Conf.Session.Name)
	if err != nil {
		logger.Errorf("<GetUserInfo> GetUserInfo error, %s", err)
		return
	}

	authToken.SetTokenInfo(models.AuthTokenInfo{
		GlobalSessionId: sessionStore.ID,
		UserName:        userInfo.UserName,
		UserId:          userInfo.UserId,
	})
	return
}
