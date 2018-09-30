package rate

import (
	"ehelp/x/db/mongodb"
)

type Rate struct {
	mongodb.BaseModel `bson:",inline"`
	CusId             string `bson:"cus_id" json:"cus_id"`
	EmpId             string `bson:"emp_id" json:"emp_id"`
	Comment           string `bson:"comment" json:"comment"`
	Rate              int   `bson:"rate" json:"rate"`
	RateBy            RateBy `bson:"rate_by" json:"rate_by"`
	OrderID           string `bson:"order_id" json:"order_id"`
}
type RateBy int

var RATE_CUS = RateBy(1)
var RATE_EMP = RateBy(2)

var RateTable = mongodb.NewTable("rate", "rate", 20)
