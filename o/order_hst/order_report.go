package order_hst

import (
	"ehelp/common"
	"ehelp/x/utils"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type OrderReport struct {
	HourDay       float64  `json:"hour_day" bson:"hour_day"`
	MoneyDay      float64  `json:"money_day" bson:"money_day"`
	CreatedAt     int      `json:"created_at" bson:"created_at"`
	Customer      string   `json:"customer" bson:"customer"`
	CustomerPhone string   `json:"customer_phone" bson:"customer_phone"`
	Employee      string   `json:"employee" bson:"employee"`
	EmployeePhone string   `json:"employee_phone" bson:"employee_phone"`
	Status        string   `json:"status" bson:"status"`
	Services      []string `json:"services" bson:"services"`
}

func join(table, localField, foreignField, as string) bson.M {
	return bson.M{
		"$lookup": bson.M{
			"from":         table,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	}
}
func GetOrderREport() ([]*OrderReport, error) {
	var joinCustomer = join("customer", "cus_id", "_id", "customer")
	var joinEmployee = join("employee", "emp_id", "_id", "employee")
	var unwindCustomer = bson.M{"$unwind": bson.M{
		"path": "$customer",
		"preserveNullAndEmptyArrays": true,
	}}
	var unwindEmployee = bson.M{"$unwind": bson.M{
		"path": "$employee",
		"preserveNullAndEmptyArrays": true,
	}}
	var project = bson.M{
		"$project": bson.M{
			"created_at":     "$created_at",
			"customer":       "$customer.full_name",
			"customer_phone": "$customer.phone",
			"employee":       "$employee.full_name",
			"employee_phone": "$employee.phone",
			"hour_day":       "$hour_day",
			"money_day":      "$money_day",
			"status":         "$item_status",
			"services":       "$services",
		},
	}
	var sort = bson.M{"$sort": bson.M{"created_at": -1}}
	var reports []*OrderReport
	err := OrderHstTable.Pipe([]bson.M{
		joinCustomer,
		joinEmployee,
		unwindCustomer,
		unwindEmployee,
		project,
		sort,
	}).All(&reports)
	return reports, err
}

type GeneralReport struct {
	Service   string  `json:"service" bson:"service"`
	Year      int     `json:"year" bson:"year"`
	Month     int     `json:"month" bson:"month"`
	Week      int     `json:"week" bson:"week"`
	Time      string  `json:"time" bson:"time"`
	Day       int     `json:"day" bson:"day"`
	CreatedAt int     `json:"created_at" bson:"created_at"`
	Total     float64 `json:"total" bson:"total"`
}

func matchGeneralReportByTime(types string) (bson.M, bson.M) {
	var firstMatch = bson.M{}
	var firstGroup = bson.M{"total": bson.M{"$sum": "$money_day"}}
	// var timeGroup = bson.M{"service": "$services"}
	var timeGroup = bson.M{}
	var now = utils.Now{
		Time: common.GetTimeNowVietNam(),
	}
	if types == "week" {
		var start, end = now.GetCurrentDay()
		firstMatch["day"] = bson.M{"$lte": end, "$gte": start}
		timeGroup["day"] = bson.M{"$dayOfMonth": "$daystr"}
	} else if types == "month" {
		var start, end = now.GetCurrentWeeks()
		firstMatch["week"] = bson.M{"$lte": end, "$gte": start}
		timeGroup["week"] = bson.M{"$week": "$daystr"}
	} else {
		var start, end = now.GetCurrentMonth()
		firstMatch["month"] = bson.M{"$lte": end, "$gte": start}

	}
	timeGroup["month"] = bson.M{"$month": "$daystr"}
	timeGroup["year"] = bson.M{"$year": "$daystr"}
	timeGroup["service"] = "$services"
	firstGroup["_id"] = timeGroup
	firstGroup["created_at"] = bson.M{"$first": "$created_at"}
	var match = bson.M{"$match": firstMatch}
	var group = bson.M{"$group": firstGroup}
	return group, match
}
func GetGeneralReportByTime(types string) ([]*GeneralReport, error) {
	var firstMatch = bson.M{"$match": bson.M{"item_status": "finished"}}
	var unwind = bson.M{"$unwind": "$services"}
	var res []*GeneralReport
	var group, match = matchGeneralReportByTime(types)
	// var unwind = bson.M{"$unwind": "$services"}
	var project = bson.M{
		"$project": bson.M{
			"service":    "$_id.service",
			"year":       "$_id.year",
			"month":      "$_id.month",
			"week":       "$_id.week",
			"day":        "$_id.day",
			"total":      "$total",
			"created_at": "$created_at",
			"_id":        0,
		},
	}
	var sort = bson.M{"$sort": bson.M{"created_at": 1}}
	err := OrderHstTable.Pipe([]bson.M{firstMatch, unwind, group, project, match, sort}).All(&res)
	if res != nil {
		for _, item := range res {
			item.TransformTime(types)
		}
	}
	return res, err
}

func (g *GeneralReport) TransformTime(types string) {
	if types == "year" {
		g.Time = strconv.Itoa(g.Month) + "/" + strconv.Itoa(g.Year)
	} else if types == "month" {
		g.Time = strconv.Itoa(g.Week)
	} else {
		g.Time = strconv.Itoa(g.Day) + "/" + strconv.Itoa(g.Month) + "/" + strconv.Itoa(g.Year)

	}
}

type Statistic struct {
	Day     int    `json:"day" bson:"day"`
	Month   int    `json:"month" bson:"month"`
	Year    int    `json:"year" bson:"year"`
	Service string `json:"service"`
	Total   int    `json:"total"`
}

func GetStatisticOrder() (ords []*Statistic, err error) {
	var query = []bson.M{}
	var groupSer = bson.M{}
	var queryMatch = bson.M{"status_item": common.ITEM_ORDER_STATUS_WORKING}
	var unWindSer = bson.M{"path": "$services", "preserveNullAndEmptyArrays": true}
	groupSer["_id"] = bson.M{
		"service": "$services",
		"year":    bson.M{"$year": "$daystr"},
		"month":   bson.M{"$month": "$daystr"},
		"day":     bson.M{"$dayOfMonth": "$daystr"},
	}
	groupSer["total"] = bson.M{"$sum": "$money_day"}
	var project = bson.M{
		"service": "$_id.service",
		"year":    "$_id.year",
		"month":   "$_id.month",
		"day":     "$_id.day",
		"total":   "$total",
		"_id":     0,
	}
	query = []bson.M{
		{"$match": queryMatch},
		{"$unwind": unWindSer},
		{"$group": groupSer},
		{"$project": project},
	}
	err = OrderHstTable.Pipe(query).All(&ords)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	}
	return ords, err
}
