package service

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"gorm.io/gorm"
	"main/pkg"
)

type Srv struct {
	repo     Repository
	log      Log
	keycloak *keycloak
}

func GetService(r Repository, l Log, conf *pkg.Config) Service {
	return &Srv{repo: r, log: l, keycloak: newKeycloak(conf)}
}

func (s *Srv) CreateUse1r() {
	s.repo.CreateUser()
}

func (s *Srv) Login(req *pkg.LoginRequest) (*pkg.LoginResponse, error) {
	jwt, err := s.kcLogin(req)
	if err != nil {
		return nil, err
	}
	s.checkUserInDb(jwt.AccessToken)
}

func (s *Srv) kcLogin(req *pkg.LoginRequest) (*gocloak.JWT, error) {
	return s.keycloak.gocloak.Login(context.Background(),
		s.keycloak.clientId,
		s.keycloak.clientSecret,
		s.keycloak.realm,
		req.Login,
		req.Password)
}

func (s *Srv) checkUserInDb(accessToken string) {
	userInfo, err := s.keycloak.gocloak.GetUserInfo(context.Background(),
		accessToken,
		s.keycloak.realm)
	user, err := s.repo.GetUserByKcId(*userInfo.Sub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

}
