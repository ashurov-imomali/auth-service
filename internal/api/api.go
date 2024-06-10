package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/internal/service"
	"main/pkg"
	"net/http"
)

type api struct {
	engine *gin.Engine
	srv    service.Service
	log    pkg.Log
}

func NewApi(s service.Service, l pkg.Log) Api {
	return &api{engine: gin.Default(), srv: s, log: l}
}

func (a *api) InitRoutes(conf *pkg.Config) {
	r := a.engine
	r.GET("/ping", a.pong)
	r.POST("/login", a.login)
	r.POST("/send-otp", a.sendOtp)
	r.POST("/confirm-otp", a.confirmOtp)
	r.GET("/refresh-token", a.refreshToken)
	gr := r.Group("/gauth")
	gr.Use(a.checkToken())
	gr.GET("/setup", a.setupGauth)
	gr.POST("/verify", a.verifyGauth)
	r.POST("/auth", a.auth)
	r.Run(fmt.Sprintf("%s:%s", conf.Srv.Host, conf.Srv.Port))
}

func (a *api) pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (a *api) login(c *gin.Context) {
	var lData pkg.LoginRequest
	if err := c.ShouldBindJSON(&lData); err != nil {
		a.log.Error(err, "couldn't parse 2 struct")
		c.Status(http.StatusBadRequest)
		return
	}
	response, err := a.srv.Login(&lData)
	if err != nil {
		a.log.Error(err, "couldn't login")
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (a *api) auth(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")
	user, err := a.srv.Auth(accessToken)
	if err != nil {
		a.log.Error(err, "couldn't auth")
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *api) refreshToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	response, err := a.srv.RefreshToken(token)
	if err != nil {
		a.log.Error(err, "couldn't get new tokens")
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (a *api) sendOtp(c *gin.Context) {
	var req pkg.OtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while binding": err.Error()})
		return
	}
	resp, hErr := a.srv.SendOTP(&req)
	if hErr != nil {
		a.log.Error(hErr.Err, hErr.Message)
		c.JSON(hErr.Status, gin.H{"message": hErr.Message})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *api) confirmOtp(c *gin.Context) {
	var otp pkg.Confirm
	if err := c.BindJSON(&otp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse 2 struct"})
		return
	}
	resp, hErr := a.srv.ConfirmOtp(&otp)
	if hErr != nil {
		a.log.Error(hErr.Err, hErr.Message)
		c.JSON(hErr.Status, gin.H{"message": hErr.Message})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *api) setupGauth(c *gin.Context) {
	userInfo, err := a.getUserInfoFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	url, hErr := a.srv.SetupGauth(userInfo.UserId, userInfo.Username)
	if hErr != nil {
		a.log.Error(hErr.Err, hErr.Message)
		c.JSON(hErr.Status, gin.H{"message": hErr.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (a *api) verifyGauth(c *gin.Context) {
	var req pkg.Confirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse 2 struct"})
		return
	}
	userinfo, err := a.getUserInfoFromContext(c)
	if err != nil {
		a.log.Error(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if hErr := a.srv.VerifyGauth(req.Value, userinfo.UserId); hErr != nil {
		a.log.Error(hErr.Err, hErr.Message)
		c.JSON(hErr.Status, gin.H{"message": hErr.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (a *api) checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		userInfo, err := a.srv.Auth(token)
		if err != nil {
			a.log.Error(err, "invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		c.Set("user_id", userInfo.UserId)
		c.Set("username", userInfo.Username)
		c.Next()
	}
}

func (a *api) getUserInfoFromContext(c *gin.Context) (*pkg.UserInfo, error) {
	userId, find := c.Get("user_id")
	if !find {
		return nil, errors.New("missing user_id")
	}
	username, find := c.Get("username")
	if !find {
		return nil, errors.New("missing username")
	}
	id, ok := userId.(int64)
	if !ok {
		return nil, errors.New(fmt.Sprintf("couldn't parse 2 int64:%v", userId))
	}
	name, ok := username.(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("couldn't parse 2 string:%v", userId))
	}
	return &pkg.UserInfo{
		UserId:   id,
		Username: name,
	}, nil
}
