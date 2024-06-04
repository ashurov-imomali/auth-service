package service

import (
	"github.com/golang-jwt/jwt/v5"
	"main/pkg"
)

type Log interface {
	Info(string)
	Error(error, string)
	Warn(string)
	Debug(string)
}

type Repository interface {
	CreateUserWithBaseRole(user *pkg.User) error
	GetUserByKcId(kcId string) (*pkg.User, error)
}

type Service interface {
	Login(data *pkg.LoginRequest) (*pkg.LoginResponse, error)
	Auth(accessToken string) (*jwt.MapClaims, error)
}
