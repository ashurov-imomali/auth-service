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
	UpdateUser(user *pkg.User) (*pkg.User, error)
	UpdateUserGauth(id int64, secret string) error
}

type Service interface {
	Login(data *pkg.LoginRequest) (*pkg.LoginResponse, *Error)
	Auth(accessToken string) (*pkg.UserInfo, *Error)
	RefreshToken(refreshToken string) (*pkg.Tokens, *Error)
	SendOTP(request *pkg.OtpRequest) (*pkg.OtpRequest, *Error)
	SetupGauth(userId int64, username string) (string, *Error)
	ConfirmOtp(otp *pkg.Confirm) (*pkg.ConfirmResp, *Error)
	VerifyGauth(otp string, userId int64) *Error
}
