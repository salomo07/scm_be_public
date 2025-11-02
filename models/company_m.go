package models

type Company struct {
	IdCompany       string    `json:"_id"`
	AppId           string    `json:"appid" validate:"required"`
	Name            string    `json:"name" validate:"required"`
	Alias           string    `json:"alias" validate:"required"`
	LevelMembership string    `json:"levelmembership"`
	Cred            string    `json:"cred"`
	Nonce           string    `json:"nonce"`
	Contacts        []Contact `json:"contacts"`
	Roles           []Role    `json:"roles"`
	Owner           User      `json:"owner" validate:"required"`
}
type Role struct {
	Id         string       `json:"_id"`
	Name       string       `json:"name" validate:"required"`
	Type       string       `json:"type" validate:"required"`
	Desc       string       `json:"desc"`
	AccessMenu []AccessMenu `json:"accessmenu"`
}
type Contact struct {
	Email    string `json:"email"` //Harus encrypt  di DB
	Phone    string `json:"phone"`
	Mobile   string `json:"mobile" validate:"required"` //Harus encrypt  di DB
	WhatsApp string `json:"whatsapp"`
}
type Branch struct {
	Id            string `bson:"_id"`
	Name          string `bson:"name"`
	Address       string `bson:"address"`
	Desc          string `bson:"desc"`
	ContactMobile string `bson:"contactmobile"`
}
type BranchRequest struct {
	Id            string `json:"_id,omitempty"`
	Name          string `json:"name" validate:"required"`
	IdCompany     string `json:"idcompany" validate:"required"`
	Address       string `json:"address" validate:"required"`
	Desc          string `json:"desc"`
	ContactMobile string `json:"contactmobile"`
}
type RoleRequest struct {
	Id         string       `json:"_id,omitempty"`
	Name       string       `json:"name" validate:"required"`
	IdCompany  string       `json:"idcompany" validate:"required"`
	Type       string       `json:"type" validate:"required"`
	Desc       string       `json:"desc"`
	AccessMenu []AccessMenu `json:"accessmenu"`
}

// type UserRequest struct {
// 	Id        string  `json:"_id,omitempty"`
// 	Name      string  `json:"name" validate:"required"`
// 	UserName  string  `json:"username" validate:"required"`
// 	IdCompany string  `json:"idcompany" validate:"required"`
// 	IdRole    string  `json:"idrole" validate:"required"`
// 	Contact   Contact `json:"contact"`
// 	AppId     string  `json:"appid"`
// 	IdBranch  string  `json:"idbranch"`
// }

type UserRequest struct {
	Id        string  `json:"_id,omitempty"`
	Appid     string  `json:"appid"`
	Contact   Contact `json:"contact"`
	Idbranch  string  `json:"idbranch"`
	Idcompany string  `json:"idcompany"`
	Idrole    string  `json:"idrole"`
	Name      string  `json:"name"`
	Nickname  string  `json:"nickname"`
	Nik       string  `json:"nik"`
	Password  string  `json:"password"`
	Username  string  `json:"username"`
	Isactive  bool    `json:"isactive"`
}

// type Contact struct {
// 	Email    string `json:"email"`
// 	Mobile   string `json:"mobile"`
// 	Phone    string `json:"phone"`
// 	Whatsapp string `json:"whatsapp"`
// }
