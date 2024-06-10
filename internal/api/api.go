package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/internal/service"
	"main/pkg"
	"net/http"
	"strconv"
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
	r.GET("/setup-gauth", a.setupGauth)
	r.POST("/verify-gauth", a.verifyGauth)
	r.GET("/refresh-token", a.refreshToken)
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
	strId := c.Query("user_id")
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing username"})
		return
	}
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse userid 2 int64"})
		return
	}

	url, hErr := a.srv.SetupGauth(id, username)
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
	strId := c.Query("user_id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong user id"})
		return
	}
	if hErr := a.srv.VerifyGauth(req.Value, id); hErr != nil {
		a.log.Error(hErr.Err, hErr.Message)
		c.JSON(hErr.Status, gin.H{"message": hErr.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
