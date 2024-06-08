package service

import (
	"main/pkg"
)

type Repository interface {
	GetUserByKcId(kcId string) (*pkg.User, bool, error)
	CreateUserWithBaseRole(user *pkg.User) error
	GetPermissionsByUserId(id int64) ([]int64, error)
	GetUserInfoByKcId(kcId string) (*pkg.UserInfo, error)
}

type Service interface {
	Login(data *pkg.LoginRequest) (*pkg.LoginResponse, error)
	Auth(accessToken string) (*pkg.UserInfo, error)
	RefreshToken(refreshToken string) (*pkg.Tokens, error)
}
