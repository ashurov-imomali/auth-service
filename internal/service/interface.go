package service

import (
	"main/pkg"
)

type Log interface {
	Info(string)
	Error(error, string)
	Warn(string)
	Debug(string)
}

type Repository interface {
	GetUserByKcId(kcId string) (*pkg.User, bool, error)
	CreateUserWithBaseRole(user *pkg.User) error
	GetPermissionsByUserId(id int64) ([]int64, error)
}

type Service interface {
	Login(data *pkg.LoginRequest) (*pkg.LoginResponse, error)
	Auth(accessToken string) (*pkg.User, error)
	RefreshToken(refreshToken string) (*pkg.LoginResponse, error)
}
