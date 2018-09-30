package employee

import (
	"ehelp/o/user"
	"ehelp/x/db/mongodb"
)

type Employee struct {
	user.UserInterface `bson:",inline"`
	Email              string       `bson:"email" json:"email"`
	EmployeeType       string       `bson:"employee_type" json:"employee_type"`
	Avatar             string       `bson:"avatar" json:"avatar"`
	Certificate        string       `bson:"certificate" json:"certificate"`
	EmployeeWork       EmployeeWork `bson:"emp_work" json:"emp_work"`
	AllCustomer        int          `bson:"all_customer" json:"all_customer"`
	AllHourWork        float32      `bson:"all_hour_work" json:"all_hour_work"`
	LinkAvatar         string       `bson:"link_avatar" json:"link_avatar"`
}

type EmployeeWork struct {
	AddressWork   string   `bson:"address_work" json:"address_work"`
	StartTimeWork float32  `bson:"start_time" json:"start_time"` // thời gian bắt đầu làm việc
	EndTimeWork   float32  `bson:"end_time" json:"end_time"`     // thời gian kết thúc
	ServiceIds    []string `bson:"service_ids" json:"service_ids"`
}

var EmployeeTable = mongodb.NewTable("employee", "emp", 12)
