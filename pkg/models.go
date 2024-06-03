package pkg

import (
	"database/sql"
	"time"
)

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

type Config struct {
	Db       *Database `json:"database"`
	KeyCloak *KeyCloak `json:"key_cloak"`
	Srv      *Server   `json:"server"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Id   int    `json:"user_id" gorm:"column:user_id"` // Уникальный идентификатор пользователя в bd
	KcId string `json:"kc_id" gorm:"column:kc_id"`     // Уникальный идентификатор пользователя в keycloak
}

func (User) TableName() string {
	return "tusers"
}

type TUser struct {
	UserID             int            `db:"user_id"`              // Уникальный идентификатор пользователя
	Username           string         `db:"username"`             // Имя пользователя
	Password           string         `db:"password"`             // Хэш пароля
	Fullname           sql.NullString `db:"fullname"`             // Полное имя пользователя
	UserDesc           sql.NullString `db:"user_desc"`            // Описание пользователя
	Email              string         `db:"email"`                // Электронная почта пользователя
	CreatedAt          time.Time      `db:"created_at"`           // Дата и время создания записи
	UpdatedAt          time.Time      `db:"updated_at"`           // Дата и время последнего обновления записи
	LoginAt            *time.Time     `db:"login_at"`             // Дата и время последнего входа
	Phone              string         `db:"phone"`                // Номер телефона пользователя
	Salt               string         `db:"salt"`                 // Соль для хэширования пароля
	CreatedBy          *int           `db:"created_by"`           // Кто создал запись
	ModifiedBy         *int           `db:"modified_by"`          // Кто последним изменил запись
	Disabled           bool           `db:"disabled"`             // Статус блокировки пользователя
	GAuthSecret        *bool          `db:"gauth_secret"`         // Секрет для двухфакторной аутентификации (Google Authenticator)
	GAuthVerified      bool           `db:"gauth_verified"`       // Флаг, указывающий, подтверждена ли двухфакторная аутентификация
	PasswordLastChange *time.Time     `db:"password_last_change"` // Дата и время последнего изменения пароля
}
