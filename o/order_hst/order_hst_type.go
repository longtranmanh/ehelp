package order_hst

import (
	"ehelp/common"
	"ehelp/x/db/mongodb"
	"time"
)

var OrderHstTable = mongodb.NewTable("order_history", "ohst", 20)

type OrderHST struct {
	mongodb.BaseModel `bson:",inline"`
	EmpId             string                 `bson:"emp_id" json:"emp_id"`
	Services          []string               `bson:"services" json:"services"`
	CusId             string                 `bson:"cus_id" json:"cus_id"`
	ItemId            string                 `bson:"item_id" json:"item_id"`
	ItemStatus        common.ItemOrderStatus `bson:"item_status" json:"item_status"`
	OrderId           string                 `bson:"order_id" json:"order_id"`
	OrderStatus       common.OrderStatus     `bson:"order_status" json:"order_status"`
	StartWorkDay      float32                `bson:"start_work_day" json:"start_work_day"`
	EndWorkDay        float32                `bson:"end_work_day" json:"end_work_day"`
	HourDay           float32                `bson:"hour_day" json:"hour_day"`
	MoneyDay          float32                `bson:"money_day" json:"money_day"`
	DayStr            time.Time              `bson:"daystr" json:"daystr"`
}
