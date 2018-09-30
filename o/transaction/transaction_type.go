package transaction

import (
	"ehelp/x/db/mongodb"
)

type Transaction struct {
	mongodb.BaseModel `bson:",inline"`
	OrderId           string  `bson:"order_id" json:"order_id"`
	UserId            string  `bson:"user_id" json:"user_id"`
	HourWork          float32 `bson:"hour_work" json:"hour_work"`
}

var TransactionTable = mongodb.NewTable("transaction", "tbx", 12)
