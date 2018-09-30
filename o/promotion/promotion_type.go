package promotion

import (
	"ehelp/x/db/mongodb"
)

type Promotion struct {
	mongodb.BaseModel `bson:",inline"`
	Content           string `bson:"content" json:"content" validate:"required"`
	Discount          struct {
		Type  string `bson:"type" json:"type" validate:"required"`
		Value int    `bson:"value" json:"value" validate:"required"`
	} `bson:"discount" json:"discount"`
	NumberOfOrder int     `bson:"number_of_order" json:"number_of_order"`
	HourStart     float32 `bson:"use_start" json:"use_start"`
	HourEnd       float32 `bson:"use_end" json:"use_end"`
	Description   string  `bson:"description" json:"description" validate:"required"`
	UrlWeb        string  `bson:"url_web" json:"url_web"`
	Title         string  `bson:"title" json:"title" validate:"required"`
}
type PromotionHistory struct {
	mongodb.BaseModel `bson:",inline"`
	PromotionID       string `bson:"promotion_id" json:"promotion_id"`
	IsActive          bool   `bson:"is_active" json:"is_active"`
	ShopID            string `bson:"customer_id" json:"customer_id"`
	OrderId           string `bson:"order_id" json:"order_id"`
}

type PromotionAndHst struct {
	Promotion `bson:",inline"`
	PrmsHst   []PromotionHistory `bson:"prms_hst" json:"prms_hst"`
}

var PromotionTable = mongodb.NewTable("promotion", "prm", 12)
var PromotionHistoryTable = mongodb.NewTable("promotion_hst", "prm_hst", 12)
