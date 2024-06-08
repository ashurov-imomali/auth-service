package service

import (
	"main/pkg"
)

type Repository interface {
	GetUserByKcId(kcId string) (*pkg.User, bool, error)
	CreateUserWithBaseRole(user *pkg.User) error
	GetPermissionsByUserId(id int64) ([]int64, error)
	GetUserInfoByKcId(kcId string) (*pkg.UserInfo, error)
	GetUserById(id int64) (*pkg.User, error)
}

type Service interface {
	Login(data *pkg.LoginRequest) (*pkg.LoginResponse, error)
	Auth(accessToken string) (*pkg.UserInfo, error)
	RefreshToken(refreshToken string) (*pkg.Tokens, error)
	SendOTP(request *pkg.OtpRequest) (*pkg.OtpRequest, *Error)
	ConfirmOtp(otp *pkg.Confirm) (*pkg.ConfirmResp, *Error)
}
