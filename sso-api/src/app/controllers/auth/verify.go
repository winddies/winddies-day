package auth

import (
	"login/src/app/code"
	"login/src/app/models"
	"login/src/app/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func (auth *auth) VerifyLogin(ctx *gin.Context) {
	logger := auth.Logger(ctx)
	session := sessions.Default(ctx)
	value := session.Get(authUserKey)
	logger.Infoln(value)
	// value != nil 说明本地浏览器有 sso_id, 则根据 session 里的 userName 去拿用户信息
	// 拿到用户信息后构造 token, 并根据 token 存储 tokenInfo，当子系统验证的时候，取出 tokenInfo 给子系统
	if value != nil {
		user := &models.User{
			UserName: value.(string),
		}
		token, err := service.HandleToken(ctx, logger, user)
		if err != nil {
			logger.Errorf("<GenerateToken> GenerateToken error, %s", err)
			auth.Send(ctx, code.ResultError, nil)
			return
		}
		auth.Send(ctx, code.OK, map[string]string{"token": token})
		return
	}

	auth.Send(ctx, code.LoginFailed, nil)
}

func (auth *auth) VerifyToken(ctx *gin.Context) {
	logger := auth.Logger(ctx)

	tokenId := ctx.Query("tokenId")
	clientSessionId := ctx.Query("localSessionId")
	clientName := ctx.Query("clientName")
	logoutUrl := ctx.Query("logoutUrl")

	if tokenId == "" {
		logger.Error("<Verify>verify failed，token is necessary")
		auth.Send(ctx, code.ResultError, nil)
		return
	}

	authToken := models.NewAuthToken(tokenId)
	info, err := authToken.GetTokenInfo()
	logger.Infof("%+v", info)
	if err != nil {
		logger.Errorf("<GetTokenInfo> authToken get tokenInfo error, %s", err)
		auth.Send(ctx, code.ResultError, nil)
		return
	}

	if info != nil {
		var tokenStructInfo models.AuthTokenInfo
		if err := mapstructure.Decode(info, &tokenStructInfo); err != nil {
			logger.Errorf("<mapstructure.Decode> decode error, %s", err)
			return
		}

		auth.Send(ctx, code.OK, tokenStructInfo)

		GlobalSessionManage := models.NewGlobalSessionManage()
		GlobalSessionManage.Set(tokenStructInfo.GlobalSessionId, &models.GlobalSessionManageInfo{
			ClientId:   clientSessionId,
			ClientName: clientName,
			LogoutUrl:  logoutUrl,
		})
	} else {
		auth.Send(ctx, code.ResultError, nil)
	}
}
