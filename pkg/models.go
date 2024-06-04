package pkg

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
	Id       int64  `json:"user_id" gorm:"column:user_id;primary_key"` // Уникальный идентификатор пользователя в bd
	KcId     string `json:"kc_id" gorm:"column:kc_id"`                 // Уникальный идентификатор пользователя в keycloak
	UserName string `json:"user_name" gorm:"column:username"`
}

func (User) TableName() string {
	return "tusers"
}

type User2Role struct {
	RoleId int64 `json:"role_id" gorm:"role_id"`
	UserId int64 `json:"user_id" gorm:"user_id"`
}

func (User2Role) TableName() string {
	return "tuser2role"
}
