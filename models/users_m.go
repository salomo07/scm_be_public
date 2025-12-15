package models

type LoginFromDB struct {
	Id        string  `json:"_id"`
	AppId     string  `json:"appid"`
	Company   Company `json:"company"`
	Contact   Contact `json:"contact"`
	IdCompany string  `json:"idcompany"`
	IdBranch  string  `json:"idbranch"`
	IdRole    string  `json:"idrole"`
	Name      string  `json:"name"`
	Nickname  string  `json:"nickname"`
	NIK       string  `json:"nik"`
	Password  string  `json:"password"`
	Role      Role    `json:"role"`
	Username  string  `json:"username"`
	Pin       string  `json:"pin"`
	Menus     []Menu  `json:"menus"`
}
type LoginResponseJWT struct {
	Id        string  `json:"_id"`
	AppId     string  `json:"appid"`
	Contact   Contact `json:"contact"`
	IdCompany string  `json:"idcompany"`
	IdBranch  string  `json:"idbranch"`
	IdRole    string  `json:"idrole"`
	Name      string  `json:"name"`
	Nickname  string  `json:"nickname"`
	NIK       string  `json:"nik"`
	Username  string  `json:"username"`
	RoleName  string  `json:"rolename"`
	RoleType  string  `json:"roletype"`
}
type User struct {
	Id        string  `json:"_id" validate:"required"`
	AppId     string  `json:"appid"`
	Nik       string  `json:"nik" validate:"required"` //Harus encrypt di DB
	Name      string  `json:"name" validate:"required"`
	Nickname  string  `json:"nickname"`
	Username  string  `json:"username" validate:"required"`
	Password  string  `json:"password" validate:"required"` //Harus encrypt  di DB
	IdCompany string  `json:"idcompany"`
	IdBranch  string  `json:"idbranch"`
	IdRole    string  `json:"idrole"`
	Pin       string  `json:"pin"`
	RoleName  string  `json:"rolename"`
	Fullname  string  `json:"fullname"`
	Contact   Contact `json:"contact" validate:"required"`
}
type UserCompany struct {
	Company Company `json:"company"`
}

type PublishRedis struct {
	IdCompany string `json:"idcompany" validate:"required"`
	Data      any    `json:"data" validate:"required"`
}
