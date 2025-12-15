package models

type AppInfo struct {
	ID      string  `json:"_id"`
	Owner   string  `json:"owner"`
	Time    float32 `json:"time"`
	Version string  `json:"version"`
	Name    string  `json:"name"`
	Code    string  `json:"code"`
}

type Menu struct {
	Id      string    `json:"_id" validate:"required"`
	Name    string    `json:"name" validate:"required"`
	Url     string    `json:"url" validate:"required"`
	Icon    string    `json:"icon"`
	Desc    string    `json:"desc"`
	Submenu []Submenu `json:"submenu"`
}
type Submenu struct {
	Id   string `json:"_id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Url  string `json:"url"`
	Icon string `json:"icon"`
	Desc string `json:"desc"`
}

type AccessMenu struct {
	Idmenu        string          `json:"idmenu" validate:"required"`
	AccessSubmenu []AccessSubmenu `json:"accesssubmenu" validate:"dive"`
	Create        *bool           `json:"create" validate:"required"`
	Read          *bool           `json:"read" validate:"required"`
	Update        *bool           `json:"update" validate:"required"`
	Delete        *bool           `json:"delete" validate:"required"`
}
type AccessSubmenu struct {
	Idsubmenu *string `json:"idsubmenu" validate:"required"`
	Create    *bool   `json:"create" validate:"required"`
	Read      *bool   `json:"read" validate:"required"`
	Update    *bool   `json:"update" validate:"required"`
	Delete    *bool   `json:"delete" validate:"required"`
}
