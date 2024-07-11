package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"io"
	"main/pkg"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	service    = "Humo ASO"
	secretSize = 15
	period     = 30
)

type Srv struct {
	repo     Repository
	log      pkg.Log
	keycloak *keycloak
	rClient  *redis.Client
	hClient  http.Client
}

func GetService(r Repository, l pkg.Log, conf *pkg.Config) Service {
	return &Srv{repo: r, log: l,
		keycloak: newKeycloak(conf.KeyCloak),
		rClient:  pkg.NewRedisClient(conf.Redis),
		hClient:  newHttpClient(conf.HClient),
	}
}

func (s *Srv) Login(req *pkg.LoginRequest) (*pkg.LoginResponse, *Error) {
	requestId := uuid.New().String()
	token, err := s.kcLogin(req)
	if err != nil {
		return nil, keyCloakError(err)
	}

	tData, err := json.Marshal(&pkg.Tokens{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken})
	if err != nil {
		return nil, internalServerError(err, "couldn't marshal")
	}

	if err := s.setRCache("token_"+requestId, tData, 8*time.Minute); err != nil {
		return nil, internalServerError(err, "redis error")
	}

	user, err := s.checkUserInDb(token.AccessToken)
	if err != nil {
		return nil, keyCloakError(err)
	}

	uData, err := json.Marshal(&pkg.UserSecure{UserID: user.Id, GauthVerified: user.GauthVerified,
		Gattribute: user.GauthSecret, Username: user.Username})
	if err != nil {
		return nil, internalServerError(err, "marshal error")
	}

	if err := s.setRCache("user_"+requestId, uData, 7*time.Minute); err != nil {
		return nil, internalServerError(err, "redis error[set]")
	}

	var gAuthSession string
	if pkg.Params.Sms2Fa {
		gAuthSession = uuid.New().String()
		if err := s.setRCache(gAuthSession, uData, 7*time.Minute); err != nil {
			return nil, internalServerError(err, "redis error[set]")
		}
	}
	permissions, err := s.repo.GetPermissionsByUserId(user.Id)
	if err != nil {
		return nil, internalServerError(err, "redis error[set]")
	}

	return &pkg.LoginResponse{
		RequestID:       requestId,
		Phone:           s.customizePhone(user.Phone),
		IsGauthPrefered: user.GauthVerified,
		SmsOtpDisable:   pkg.Params.Sms2Fa,
		GauthSession:    gAuthSession,
		FirstLogin:      len(permissions) == 0,
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
		user := &pkg.User{KcId: *userInfo.Sub,
			Username: *userInfo.PreferredUsername,
			Disabled: true,
			FullName: checkValue(userInfo.Name, "").(string),
			Email:    checkValue(userInfo.Email, "").(string),
		}
		return user, s.repo.CreateUserWithBaseRole(user)
	}
	return user, err
}

func (s *Srv) getKcUserInfo(accessToken string) (*gocloak.UserInfo, error) {
	return s.keycloak.gocloak.GetUserInfo(context.Background(),
		accessToken,
		s.keycloak.realm)
}

func (s *Srv) Auth(accessToken string) (*pkg.UserInfo, *Error) {
	result, err := s.verifyToken(accessToken)
	if err != nil {
		return nil, keyCloakError(err)
	}
	if !*result.Active {
		return nil, unauthorized(errors.New("invalid token"), "invalid token")
	}

	_, claims, err := s.decodeAccessToken(accessToken)
	if err != nil {
		return nil, keyCloakError(err)
	}

	user, err := s.repo.GetUserInfoByKcId((*claims)["sub"].(string))
	if err != nil {
		return nil, internalServerError(err, "database error")
	}
	return user, nil
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

func (s *Srv) RefreshToken(refreshToken string) (*pkg.Tokens, *Error) {
	token, err := s.refreshToken(refreshToken)
	if err != nil {
		return nil, keyCloakError(err)
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

func (s *Srv) SendOTP(req *pkg.OtpRequest) (*pkg.OtpRequest, *Error) {
	var usrSecrets pkg.UserSecure
	redisNil, err := s.getRCache("user_"+req.RequestID, &usrSecrets)
	if redisNil {
		return nil, unauthorized(err, "First you have login")
	}
	if err != nil {
		return nil, internalServerError(err, "Redis Error")
	}
	user, err := s.repo.GetUserById(usrSecrets.UserID)
	if err != nil {
		return nil, internalServerError(err, "Database error")
	}
	otpSms := pkg.SmsOTP{
		ID:      strconv.FormatInt(user.Id, 10) + req.RequestID,
		Account: strings.ReplaceAll(user.Phone, " ", ""),
	}
	otp, hErr := s.sendSmsOtp(&otpSms)
	if hErr != nil {
		return nil, hErr
	}
	usrSecrets.OtpID = otp.ID
	data, err := json.Marshal(&usrSecrets)
	if err != nil {
		return nil, internalServerError(err, "couldn't unmarshal")
	}
	if err := s.setRCache("user_"+req.RequestID, data, 8*time.Minute); err != nil {
		return nil, internalServerError(err, "Redis error")
	}
	return req, nil
}

func (s *Srv) sendSmsOtp(otp *pkg.SmsOTP) (*pkg.SmsOTP, *Error) {
	otp.Lifetime = pkg.Params.OTPLifetime
	otp.ConfirmLimit = pkg.Params.OTPConfirmLimit

	var responseOTP pkg.SmsOTP
	marshal, err := json.Marshal(otp)
	if err != nil {
		return nil, internalServerError(err, "Couldn't marshal struct")
	}
	requestBody := bytes.NewBuffer(marshal)
	req, err := http.NewRequest(http.MethodPost, pkg.Params.OTPUrl, requestBody)
	if err != nil {
		return nil, internalServerError(err, "couldn't parse 2 struct")
	}
	s.log.Info(fmt.Sprintf("request:%v", req))
	resp, err := s.hClient.Do(req)
	if err != nil {
		return nil, internalServerError(err, "Send otp unavailable")

	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, badRequest(errors.New("otp adapter error"), string(body))
	}

	err = json.Unmarshal(body, &responseOTP)
	if err != nil {
		return nil, internalServerError(err, "couldn't unmarshal response from otp")
	}
	return &responseOTP, nil
}

func (s *Srv) ConfirmOtp(otp *pkg.Confirm) (*pkg.ConfirmResp, *Error) {
	var usrSecrets pkg.UserSecure
	redisNil, err := s.getRCache("user_"+otp.RequestID, &usrSecrets)
	if redisNil {
		return nil, unauthorized(err, "first u have login")
	}
	if err != nil {
		return nil, internalServerError(err, "redis error")
	}

	switch otp.Type {
	case "sms":
		otpSms := pkg.SmsOTP{
			ID:    usrSecrets.OtpID,
			Value: otp.Value,
		}
		if hErr := s.confirmSmsOtp(&otpSms); hErr != nil {
			return nil, hErr
		}
	case "gauth":
		if !usrSecrets.GauthVerified {
			return nil, unauthorized(errors.New("wrong type"), "wrong otp type")
		}
		if validate := totp.Validate(otp.Value, usrSecrets.Gattribute); !validate {
			return nil, unauthorized(errors.New("invalid google totp"), "invalid totp")

		}
	default:
		return nil, unauthorized(errors.New("unknown OTP type"), "Unknown OTP type")
	}
	permissions, err := s.repo.GetPermissionsByUserId(usrSecrets.UserID)
	if err != nil {
		return nil, internalServerError(err, "Database error")
	}
	var token pkg.Tokens
	redisNil, err = s.getRCache("token_"+otp.RequestID, &token)
	if redisNil {
		return nil, unauthorized(err, "first u have login")
	}
	if err != nil {
		return nil, internalServerError(err, "redis error")
	}
	return &pkg.ConfirmResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserId:       usrSecrets.Username,
		Permissions:  permissions,
	}, nil
}

func (s *Srv) confirmSmsOtp(otpSms *pkg.SmsOTP) *Error {
	marshal, err := json.Marshal(otpSms)
	if err != nil {
		return badRequest(err, "couldn't parse 2 struct")
	}
	requestBody := bytes.NewBuffer(marshal)

	req, _ := http.NewRequest(http.MethodPatch, pkg.Params.OTPUrl+"/"+otpSms.ID, requestBody)
	s.log.Info(fmt.Sprintf("request:%v", req))
	resp, err := s.hClient.Do(req)
	if err != nil {
		return internalServerError(err, "couldn't send to otp")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return internalServerError(err, "otp service error")
	}
	if resp.StatusCode != http.StatusOK {
		return unauthorized(errors.New("otp adapter error"), string(body)+resp.Status)
	}

	return nil
}

func (s *Srv) SetupGauth(userId int64, username string) (string, *Error) {
	key, err := s.generateGauth(username)
	if err != nil {
		return "", internalServerError(err, "Couldn't generate gauth secret")
	}
	if err := s.repo.UpdateUserGauth(userId, key.Secret()); err != nil {
		return "", internalServerError(err, "Database error[UpdateUSER]")
	}
	return key.URL(), nil
}

func (s *Srv) generateGauth(username string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Period:      period,
		Issuer:      service,
		AccountName: username,
		SecretSize:  secretSize,
	})
}

func (s *Srv) VerifyGauth(otp string, userId int64) *Error {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return internalServerError(err, "Database error[GetUSER]")
	}
	if ok := totp.Validate(otp, user.GauthSecret); !ok {
		return unauthorized(errors.New("invalid otp"), "invalid otp")
	}
	user.GauthVerified = true
	if _, err := s.repo.UpdateUser(user); err != nil {
		return internalServerError(err, "Database error[UpdateUSER]")
	}
	return nil
}

func checkValue(v interface{}, zeroVal interface{}) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		return val.Elem().Interface()
	}
	return zeroVal
}
