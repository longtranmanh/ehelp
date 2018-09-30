package order

import (
	"ehelp/common"
	"ehelp/x/db/mongodb"
)

type Order struct {
	mongodb.BaseModel     `bson:",inline"`
	EmpID                 string    `bson:"emp_id" json:"emp_id"`
	CusID                 string    `bson:"cus_id" json:"cus_id" validate:"required"`
	AddressLoc            *Location `bson:"address_loc" json:"address_loc"`
	PricePromotion        float32   `bson:"price_promotion" json:"price_promotion"`
	PriceAllHour          float32   `bson:"price_all_hour" json:"price_all_hour"`
	PriceTool             float32   `bson:"price_tool" json:"price_tool"`
	Rated                 bool      `bson:"rated" json:"rated"`
	PriceEnd              int       `bson:"price_end" json:"price_end"`
	Note                  string    `bson:"note" json:"note"`
	common.MathPriceOrder `bson:",inline"`
	AllHourWork           float32            `bson:"all_hour_work" json:"all_hour_work" validate:"required"`
	Status                common.OrderStatus `bson:"status" json:"status" validate:"required"`
	Employee              *CustomerEmp       `bson:"employee" json:"employee"`
	Customer              *CustomerEmp       `bson:"customer" json:"customer"`
	Services              []*ServiceDetail   `bson:"services" json:"services"`
}

type ServiceDetail struct {
	ID           string   `json:"id" bson:"id"`
	Name         string   `json:"name" bson:"name"`
	NodeServices []string `json:"node_services" bson:"node_services"`
}

type Location struct {
	Address string  `bson:"address" json:"address" validate:"required"`
	Lat     float64 `bson:"lat" json:"lat" validate:"required"`
	Lng     float64 `bson:"lng" json:"lng" validate:"required"`
}

type ServiceWork struct {
	ServiceId  string `bson:"service_id" json:"service_id"`
	NumberTime int    `bson:"number_time" json:"number_time"`
}

type ToolWork struct {
	ToolId string `bson:"tool_id" json:"tool_id"`
	Number int    `bson:"num" json:"num"`
}

var OrderTable = mongodb.NewTable("order", "ord", 12)
