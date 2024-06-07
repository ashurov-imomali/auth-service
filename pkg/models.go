package pkg

import "time"

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type KeyCloak struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Realm        string `json:"realm"`
}

type RedisClient struct {
}

type Config struct {
	Db       *Database `json:"database"`
	KeyCloak *KeyCloak `json:"key_cloak"`
	Srv      *Server   `json:"server"`
}

type LoginRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	UserId       int64   `json:"user_id"`
	Permissions  []int64 `json:"permissions"`
}

type User struct {
	Id                 int64      `gorm:"column:user_id;primary_key"`
	KcId               string     `json:"kc_id" gorm:"column:kc_id"` // Уникальный идентификатор пользователя в keycloak
	Username           string     `gorm:"column:username"`
	Password           string     `gorm:"column:password"`
	Email              string     `gorm:"column:email"`
	Desc               string     `gorm:"column:user_desc"`
	FullName           string     `gorm:"column:fullname"`
	Phone              string     `gorm:"column:phone"`
	Salt               string     `gorm:"column:salt"`
	Disabled           bool       `gorm:"column:disabled"`
	CreatedAt          *time.Time `gorm:"column:created_at"`
	UpdatedAt          *time.Time `gorm:"column:updated_at"`
	LoginAt            *time.Time `gorm:"column:login_at"`
	GauthSecret        string     `gorm:"gauth_secret"`
	GauthVerified      bool       `gorm:"gauth_verified"`
	PasswordLastChange time.Time  `gorm:"column:password_last_change"`
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
