package service

import "main/pkg"

type Log interface {
	Info(string)
	Error(error, string)
	Warn(string)
	Debug(string)
}

type Repository interface {
	CreateUser()
	GetUserByKcId(kcId string) (*pkg.User, error)
}

type Service interface {
	CreateUse1r()
	Login(data *pkg.LoginRequest) (*pkg.LoginResponse, error)
}
