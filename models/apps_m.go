package models

type Menu struct {
	Id      string    `json:"_id" validate:"required"`
	AppId   string    `json:"appid" validate:"required"`
	Name    string    `json:"name" validate:"required"`
	Code    string    `json:"code" validate:"required"`
	Url     string    `json:"url" validate:"required"`
	Icon    string    `json:"icon"`
	Desc    string    `json:"desc"`
	Submenu []Submenu `json:"submenu"`
}
type Submenu struct {
	IdSubmenu int    `json:"idsubmenu" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Url       string `json:"url"`
	Icon      string `json:"icon"`
	Desc      string `json:"desc"`
}

type AccessMenu struct {
	Id            string          `json:"_id" validate:"required"`
	IdRole        string          `json:"idrole" validate:"required"`
	Idmenu        string          `json:"idmenu" validate:"required"`
	Menuname      string          `json:"menuname" validate:"required"`
	AccessSubmenu []AccessSubmenu `json:"accesssubmenu" validate:"dive"`
	Create        *bool           `json:"create" validate:"required"`
	Read          *bool           `json:"read" validate:"required"`
	Update        *bool           `json:"update" validate:"required"`
	Delete        *bool           `json:"delete" validate:"required"`
}
type AccessSubmenu struct {
	Idsubmenu *int  `json:"idsubmenu" validate:"required"`
	Create    *bool `json:"create" validate:"required"`
	Read      *bool `json:"read" validate:"required"`
	Update    *bool `json:"update" validate:"required"`
	Delete    *bool `json:"delete" validate:"required"`
}
