package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"main/pkg"
	"strconv"
	"strings"
	"time"
)

type Srv struct {
	repo     Repository
	log      pkg.Log
	keycloak *keycloak
	rClient  *redis.Client
}

func GetService(r Repository, l pkg.Log, conf *pkg.Config) Service {
	return &Srv{repo: r, log: l, keycloak: newKeycloak(conf.KeyCloak), rClient: pkg.NewRedisClient(conf.Redis)}
}

func (s *Srv) Login(req *pkg.LoginRequest) (*pkg.LoginResponse, error) {
	requestId := uuid.New().String()
	token, err := s.kcLogin(req)
	if err != nil {
		return nil, err
	}

	tData, err := json.Marshal(&pkg.Tokens{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken})
	if err != nil {
		return nil, err
	}

	if err := s.setRCache("token_"+requestId, tData, 8*time.Minute); err != nil {
		return nil, err
	}

	user, err := s.checkUserInDb(token.AccessToken)
	if err != nil {
		return nil, err
	}

	uData, err := json.Marshal(&pkg.UserSecure{UserID: strconv.FormatInt(user.Id, 10), GauthVerified: user.GauthVerified,
		Gattribute: user.GauthSecret, Username: user.Username})
	if err != nil {
		return nil, err
	}

	if err := s.setRCache("user_"+requestId, uData, 7*time.Minute); err != nil {
		return nil, err
	}

	var gAuthSession string
	if pkg.Sms2FA {
		gAuthSession = uuid.New().String()
		if err := s.setRCache(gAuthSession, uData, 7*time.Minute); err != nil {
			return nil, err
		}
	}

	return &pkg.LoginResponse{
		RequestID:       requestId,
		Phone:           s.customizePhone(user.Phone),
		IsGauthPrefered: user.GauthVerified,
		SmsOtpDisable:   pkg.Sms2FA,
		GauthSession:    gAuthSession,
	}, nil

}

func (s *Srv) customizePhone(phone string) string {
	switch len(phone) {
	case 12:
		return phone[:5] + "*****" + phone[10:]
	case 9:
		return phone[:2] + "*****" + phone[7:]
	}
	return phone
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

func (s *Srv) RefreshToken(refreshToken string) (*pkg.Tokens, error) {
	token, err := s.refreshToken(s.extractBearerToken(refreshToken))
	if err != nil {
		return nil, err
	}
	return &pkg.Tokens{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}, nil
}

func (s *Srv) refreshToken(refreshToken string) (*gocloak.JWT, error) {
	return s.keycloak.gocloak.RefreshToken(context.Background(),
		refreshToken,
		s.keycloak.clientId,
		s.keycloak.clientSecret,
		s.keycloak.realm)
}
