package auth

import (
	"login/src/app/code"
	"login/src/app/controllers/base"
	"login/src/app/global"
	"login/src/app/models"
	"login/src/app/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type auth struct {
	*base.Base
}

func New() Auth {
	return &auth{}
}

type Auth interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	VerifyToken(ctx *gin.Context)
	VerifyLogin(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type LoginOptions struct {
	*models.User
	RedirectUrl string
}

const authUserKey string = "authUser"

func (auth *auth) Login(ctx *gin.Context) {
	var loginOptions LoginOptions
	logger := auth.Logger(ctx)
	err := ctx.BindJSON(&loginOptions)
	if err != nil {
		logger.Errorf("<Login> BindJSON error, %s", err)
		auth.Send(ctx, code.LoginFailed, nil)
		return
	}

	exist, err := loginOptions.User.VerifyAuth()
	if err != nil {
		logger.Errorf("<Login> VerifyAuth error, %s", err)
		auth.Send(ctx, code.LoginFailed, nil)
		return
	}

	if !exist {
		auth.Send(ctx, code.LoginFailed, nil)
		return
	}

	session := sessions.Default(ctx)
	session.Set(authUserKey, loginOptions.UserName)
	err = session.Save()
	if err != nil {
		logger.Errorf("<Login> session.Save error, %s", err)
		auth.Send(ctx, code.LoginFailed, nil)
		return
	}

	// 重定向到指定的 url
	if loginOptions.RedirectUrl != "" {
		token, err := service.HandleToken(ctx, logger, loginOptions.User)
		if err != nil {
			logger.Errorf("<GenerateToken> GenerateToken error, %s", err)
			auth.Send(ctx, code.ResultError, nil)
			return
		}
		logger.Infoln(loginOptions.RedirectUrl + "?token=" + token)

		auth.Send(ctx, code.OK, loginOptions.RedirectUrl+"?token="+token)
	} else {
		auth.Send(ctx, code.OK, nil)
	}
}

func (auth *auth) Register(ctx *gin.Context) {
	var user models.User
	logger := auth.Logger(ctx)
	err := ctx.BindJSON(&user)
	if err != nil {
		logger.Errorf("<Register> BindJSON error, %s", err)
		auth.Send(ctx, code.ResultError, nil)
		return
	}

	exist, err := user.Register()
	if err != nil {
		logger.Error("register error")
		auth.Send(ctx, code.ResultError, nil)
		return
	}

	if exist {
		logger.Info("user is exist")
		auth.Send(ctx, code.RegisterError, nil)
		return
	}

	auth.Send(ctx, code.OK, nil)
}

func (auth *auth) Logout(ctx *gin.Context) {
	logger := auth.Logger(ctx)

	globalId := ctx.Param("globalId")
	GlobalSessionManage := models.NewGlobalSessionManage()
	// 带 globalId 说明是从客户端退出，则 sso 这边是没有自身的 session 的
	if globalId != "" {
		err := GlobalSessionManage.DeleteGlobalSession(globalId)

		if err != nil {
			logger.Errorf("<DeleteGlobalSession> GlobalSessionManage DeleteGlobalSession error, %s", err)
			return
		}

		clientIds, err := GlobalSessionManage.GetClientIds(globalId)
		if err != nil {
			logger.Errorf("<GetClientIds> GlobalSessionManage GetClientIds error, %s", err)
			return
		}
		service.NotifyClientLogout(logger, clientIds, GlobalSessionManage)
		return
	}

	// 通过 sso 自身的 logout 退出，能拿到从浏览器而来的 sessionId
	session := sessions.Default(ctx)
	session.Delete(authUserKey)
	session.Save()
	auth.Send(ctx, code.OK, nil)

	sessionStore, err := models.RedisSessionStore.Get(ctx.Request, global.Conf.Session.Name)

	if sessionStore.ID == "" {
		logger.Error("<RedisSessionStore.Get> globalId is necessary to find child system")
		return
	}

	clientIds, err := GlobalSessionManage.GetClientIds(sessionStore.ID)
	if err != nil {
		logger.Errorf("<GetClientIds> GlobalSessionManage GetClientIds error, %s", err)
		return
	}

	service.NotifyClientLogout(logger, clientIds, GlobalSessionManage)
	return
}
