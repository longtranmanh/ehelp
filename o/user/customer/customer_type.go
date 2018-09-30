package customer

import (
	"ehelp/o/user"
	"ehelp/x/db/mongodb"
)

type Customer struct {
	user.UserInterface `bson:",inline"`
	Email              string   `bson:"email" json:"email" validate:"omitempty,email"`
	CustomerType       string   `bson:"customer_type" json:"customer_type"`
	Promotions         []string `bson:"promotions" json:"promotions"`
}

var CustomerTable = mongodb.NewTable("customer", "cus", 12)
