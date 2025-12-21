package models

// LoginRequest defines the request body for user login.
// @Description The structure of login request with username and password.
type LoginRequest struct {
	// @description The company ID of the user (optional).
	IdCompany string `json:"idcompany"`

	// @description Whether to remember the user's password (optional).
	RememberPassword bool `json:"rememberpassword"`

	// @description The duration of the login session in hours (optional).
	Duration int `json:"duration"` //	Duration on hour (optional)

	// @required
	// @description The username of the user.
	Username string `json:"username" validate:"required"`

	// @description The username of the user.
	Mobile string `json:"mobile"`

	// @required
	// @description The password of the user.
	Password string `json:"password" validate:"required"`
}

// LoginResponse defines the response for user login.
// @Description The structure of the login response.
// LoginResponse defines the response for user login.
// @Description The structure of the login response.
type LoginResponse struct {
	// @description The company ID of the user (mandatory).
	IdCompany string `json:"idcompany"`

	IdUser string `json:"iduser"`

	IdRole string `json:"idrole"`

	Username string `json:"username"`

	Fullname string `json:"fullname"`

	RoleName string `json:"rolename"`

	RoleType string `json:"roletype"`

	// @description The token generated for the user after login.
	Token string `json:"token"`

	RefreshToken string `json:"refresh_token"`

	// @description The expiration time of the generated token.
	Expired string `json:"expired"`

	// @description Indicates whether the user is successfully logged in.
	IsLogin bool `json:"islogin"`

	// @description Indicates whether the system is under maintenance.
	IsMaintenance bool           `json:"ismaintenance"`
	Menus         []MenuForLogin `json:"menus"`
	AppInfo       AppInfo        `json:"app"`
	NeedShift     bool           `json:"needshiftsession"`
}
type LoginMenu struct {
	Id      string         `json:"_id"`
	Name    string         `json:"name"`
	Code    string         `json:"code"`
	Url     string         `json:"url"`
	Icon    string         `json:"icon"`
	Desc    string         `json:"desc"`
	Create  *bool          `json:"create"`
	Read    *bool          `json:"read"`
	Update  *bool          `json:"update"`
	Delete  *bool          `json:"delete"`
	Submenu []LoginSubmenu `json:"submenu"`
}

type LoginSubmenu struct {
	IdSubmenu int    `json:"idsubmenu"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	Icon      string `json:"icon"`
	Desc      string `json:"desc"`
	Create    *bool  `json:"create"`
	Read      *bool  `json:"read"`
	Update    *bool  `json:"update"`
	Delete    *bool  `json:"delete"`
}
type OTPRequest struct {
	UserId string `json:"user_id" validate:"required"`
}
type OTPResponse struct {
	UserId  string `json:"user_id"`
	OTPCode string `json:"opt"`
}
