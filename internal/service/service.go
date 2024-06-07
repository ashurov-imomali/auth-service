package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"main/pkg"
	"os"
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
	token, err := s.kcLogin(req)
	if err != nil {
		return nil, err
	}
	user, err := s.checkUserInDb(token.AccessToken)
	if err != nil {
		return nil, err
	}
	return &pkg.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserId:       user.Id,
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

func (s *Srv) checkUserInDb(accessToken string) (*pkg.User, error) {
	userInfo, err := s.getKcUserInfo(accessToken)
	if err != nil {
		return nil, err
	}
	bytes, err := json.MarshalIndent(&userInfo, "", " ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile("datafromAD.json", bytes, 0777); err != nil {
		return nil, err
	}
	user, find, err := s.repo.GetUserByKcId(*userInfo.Sub)
	if !find {
		user := &pkg.User{KcId: *userInfo.Sub, Username: *userInfo.PreferredUsername, Disabled: true}
		return user, s.repo.CreateUserWithBaseRole(user)
	}
	return user, err
}

func (s *Srv) getKcUserInfo(accessToken string) (*gocloak.UserInfo, error) {
	return s.keycloak.gocloak.GetUserInfo(context.Background(),
		accessToken,
		s.keycloak.realm)
}

func (s *Srv) Auth(accessToken string) (*pkg.UserInfo, error) {
	result, err := s.verifyToken(s.extractBearerToken(accessToken))
	if err != nil {
		return nil, err
	}
	if !*result.Active {
		return nil, errors.New("invalid token")
	}

	_, claims, err := s.decodeAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserInfoByKcId((*claims)["sub"].(string))
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

//func (s *Srv) generateClaims(kcId string) {
//	user, find, err := s.repo.GetUserByKcId(kcId)
//	if err != nil {
//
//	}
//	mp := make(map[string]interface{})
//	mp[""]
//}

func (s *Srv) RefreshToken(refreshToken string) (*pkg.LoginResponse, error) {
	token, err := s.refreshToken(s.extractBearerToken(refreshToken))
	if err != nil {
		return nil, err
	}
	return &pkg.LoginResponse{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}, nil
}

func (s *Srv) refreshToken(refreshToken string) (*gocloak.JWT, error) {
	return s.keycloak.gocloak.RefreshToken(context.Background(),
		refreshToken,
		s.keycloak.clientId,
		s.keycloak.clientSecret,
		s.keycloak.realm)
}
