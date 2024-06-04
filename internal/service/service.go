package service

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"main/pkg"
	"strings"
)

type Srv struct {
	repo     Repository
	log      Log
	keycloak *keycloak
}

func GetService(r Repository, l Log, conf *pkg.Config) Service {
	return &Srv{repo: r, log: l, keycloak: newKeycloak(conf)}
}

func (s *Srv) Login(req *pkg.LoginRequest) (*pkg.LoginResponse, error) {
	jwt, err := s.kcLogin(req)
	if err != nil {
		return nil, err
	}
	if err := s.checkUserInDb(jwt.AccessToken); err != nil {
		return nil, err
	}
	return &pkg.LoginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}, nil

}

func (s *Srv) kcLogin(req *pkg.LoginRequest) (*gocloak.JWT, error) {
	return s.keycloak.gocloak.Login(context.Background(),
		s.keycloak.clientId,
		s.keycloak.clientSecret,
		s.keycloak.realm,
		req.Login,
		req.Password)
}

func (s *Srv) checkUserInDb(accessToken string) error {
	userInfo, err := s.getKcUserInfo(accessToken)
	if err != nil {
		return err
	}
	user, err := s.repo.GetUserByKcId(*userInfo.Sub)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.repo.CreateUserWithBaseRole(user)
	}
	return err
}

func (s *Srv) getKcUserInfo(accessToken string) (*gocloak.UserInfo, error) {
	return s.keycloak.gocloak.GetUserInfo(context.Background(),
		accessToken,
		s.keycloak.realm)
}

func (s *Srv) Auth(accessToken string) (*jwt.MapClaims, error) {
	result, err := s.verifyToken(accessToken)
	if err != nil || !*result.Active {
		return nil, err
	}

	_, claims, err := s.decodeAccessToken(accessToken)
	return claims, err
}

func (s *Srv) decodeAccessToken(accessToken string) (*jwt.Token, *jwt.MapClaims, error) {
	return s.keycloak.gocloak.DecodeAccessToken(context.Background(),
		accessToken,
		s.keycloak.realm)
}

func (s *Srv) verifyToken(accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	return s.keycloak.gocloak.RetrospectToken(context.Background(),
		accessToken,
		s.keycloak.clientId,
		s.keycloak.clientSecret,
		s.keycloak.realm)
}

func (s *Srv) extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}
