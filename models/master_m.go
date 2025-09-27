package models

type HandlingUnits struct {
	Id                 string `json:"_id"`
	HUCode             string `json:"hu_code"`
	Seq_Num            int    `json:"seq_num" validate:"required"`
	LocationId         string `json:"location_id"`
	CreateTime         int    `json:"create_time"`
	IsActive           bool   `json:"is_active"`
	HandlingUnitTypeId string `json:"handling_unit_type_id" validate:"required"`
}
type HandlingUnitType struct {
	Id          string  `json:"_id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	LengthCm    int     `json:"length_cm"`
	WidthCm     int     `json:"width_cm"`
	HeightCm    int     `json:"height_cm"`
	MaxWeightKg float32 `json:"max_weight_kg"`
	IsActive    bool    `json:"is_active"`
}

type PruductMenu struct { //Untuk Eatrary
	Id          string   `json:"_id"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Price       int      `json:"price" validate:"required"`
	Images      []string `json:"images"`
	Category    string   `json:"category"`
	IsActive    bool     `json:"is_active"`
}

//	type HUCapacity struct {
//		Id       string  `json:"_id"`
//		IdGoods  string  `json:"id_goods" validate:"required"`
//		HUTypeId string  `json:"hu_type_id" validate:"required"`
//		Capacity float32 `json:"capacity"`
//		Desc     string  `json:"desc"`
//	}
type Goods struct {
	Id          string `json:"_id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	Category    string `json:"category"`
	IsActive    string `json:"is_active"`
}
