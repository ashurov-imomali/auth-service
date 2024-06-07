package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/internal/service"
	"main/pkg"
	"net/http"
)

type api struct {
	engine *gin.Engine
	srv    service.Service
	log    service.Log
}

func NewApi(s service.Service, l service.Log) Api {
	return &api{engine: gin.Default(), srv: s, log: l}
}

func (a *api) InitRoutes(conf *pkg.Config) {
	r := a.engine
	r.GET("/ping", a.pong)
	r.POST("/login", a.login)
	r.GET("/refreshToken", a.refreshToken)
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
