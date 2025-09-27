package models

type ImageItem struct {
	ID    string  `json:"id" bson:"id"`
	Src   string  `json:"src" bson:"src"`
	IsNew bool    `json:"isNew" bson:"isNew"`
	Order float64 `json:"order" bson:"order"` // pakai float64 kalau mau mirip number JS
}

type AttributeProduct struct {
	Type            string `json:"type" bson:"type"`
	Name            string `json:"name" bson:"name"`
	PriceAdjustment int32  `json:"priceAdjustment" bson:"priceAdjustment"`
	Description     string `json:"description" bson:"description"`
}

type ProductDelete struct {
	Id string `json:"idproduct" bson:"_id" validate:"required"`
}

// ---- interface umum ----
type WithId interface {
	GetId() string
}

// ---- request struct ----
type Product struct {
	Id          string             `json:"_id,omitempty"`
	IdBranch    string             `json:"idbranch" bson:"idbranch"`
	Name        string             `json:"name" bson:"name"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	Description string             `json:"description" bson:"description"`
	IsActive    bool               `json:"isactive" bson:"isactive"`
	Images      []ImageItem        `json:"images" bson:"images"`
	CookingTime float64            `json:"cookingtime" bson:"cookingtime"`
	Attribute   []AttributeProduct `json:"attribute" bson:"attribute"`
}

func (p Product) GetId() string { return p.Id }

type Category struct {
	Id          string `json:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"desc" bson:"desc"`
}

func (c Category) GetId() string { return c.Id }
