package pkg

import (
	"time"
)

type Database struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Dbname   string `json:"dbname" yaml:"dbname"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

type Server struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}
type Redis struct {
	Uri      string `json:"uri" yaml:"uri"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
}

type KeyCloak struct {
	Host         string `json:"host" yaml:"host"`
	Port         string `json:"port" yaml:"port"`
	ClientId     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	Realm        string `json:"realm" yaml:"realm"`
}

type RedisClient struct {
}

type Config struct {
	Db        *Database         `json:"database" yaml:"database"`
	KeyCloak  *KeyCloak         `json:"key_cloak" yaml:"key_cloak"`
	Srv       *Server           `json:"server" yaml:"server"`
	Redis     *Redis            `json:"redis" yaml:"redis"`
	TFAParams *TFAParams        `json:"2fa_params" yaml:"2fa_params"`
	HClient   *HttpClientParams `json:"h_client" yaml:"h_client"`
}

type HttpClientParams struct {
	Timeout int `json:"timeout"`
}

type TFAParams struct {
	Sms2Fa          bool   `json:"sms_2fa" yaml:"sms_2fa"`
	OTPUrl          string `json:"otp_url" yaml:"otp_url"`
	OTPLifetime     int64  `json:"otp_lifetime" yaml:"otp_lifetime"` //seconds
	OTPConfirmLimit int64  `json:"otp_confirm_limit" yaml:"otp_confirm_limit"`
}

type LoginRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	RequestID       string `json:"requestID"`
	Phone           string `json:"phone"`
	IsGauthPrefered bool   `json:"isGauthPrefered"`
	SmsOtpDisable   bool   `json:"sms_otp_disable"`
	GauthSession    string `json:"gauth_session"`
	FirstLogin      bool   `json:"first_login"`
}

type User struct {
	Id            int64      `gorm:"column:user_id;primary_key"`
	KcId          string     `json:"kc_id" gorm:"column:kc_id"` // Уникальный идентификатор пользователя в keycloak
	Username      string     `gorm:"column:username"`
	Password      string     `gorm:"column:password"`
	Email         string     `gorm:"column:email"`
	Desc          string     `gorm:"column:user_desc"`
	FullName      string     `gorm:"column:fullname"`
	Phone         string     `gorm:"column:phone"`
	Salt          string     `gorm:"column:salt"`
	Disabled      bool       `gorm:"column:disabled"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	LoginAt       *time.Time `gorm:"column:login_at"`
	GauthSecret   string     `gorm:"gauth_secret"`
	GauthVerified bool       `gorm:"gauth_verified"`
}

func (User) TableName() string {
	return "tusers"
}

type UserInfo struct {
	UserId   int64  `json:"user_id" gorm:"column:user_id;primary_key"`
	KcId     string `json:"kc_id" gorm:"column:kc_id"` // Уникальный идентификатор пользователя в keycloak
	Username string `json:"username" gorm:"column:username"`
	Role     string `json:"role" gorm:"column:role"`
}

type User2Role struct {
	RoleId int64 `json:"role_id" gorm:"role_id"`
	UserId int64 `json:"user_id" gorm:"user_id"`
}

func (User2Role) TableName() string {
	return "tuser2role"
}

type UserSecure struct {
	UserID        int64
	OtpID         string
	Username      string
	GauthVerified bool
	Gattribute    string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OtpRequest struct {
	RequestID string `json:"requestID"`
}
type SmsOTP struct {
	ID           string `json:"id"`
	Account      string `json:"account"`
	Value        string `json:"value"`
	Lifetime     int64  `json:"lifetime"`
	ConfirmLimit int64  `json:"validate_limit"`
	State        int64  `json:"state"`
	CreatedAt    string `json:"created_at"`
	ExpiredAt    string `json:"expired_at"`
}

type Confirm struct {
	RequestID string `json:"requestID"`
	Value     string `json:"value"`
	Type      string `json:"type"` //could be gauth or sms
}

type ConfirmResp struct {
	AccessToken  string  `json:"token"`
	RefreshToken string  `json:"refresh_token"`
	UserId       string  `json:"userId"`
	Permissions  []int64 `json:"permissions"`
}
