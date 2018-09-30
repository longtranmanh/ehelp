package order_hst

import (
	"ehelp/common"
	"ehelp/x/mlog"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

var logOrdHst = mlog.NewTagLog("OrderHST")

func (ordHst *OrderHST) CrateOrderHistory() (*OrderHST, error) {
	ordHst.DayStr = common.GetTimeNowVietNam()
	err := OrderHstTable.Create(ordHst)
	if err != nil {
		logOrdHst.Errorf("CrateOrderHistory", err)
	}
	return ordHst, err
}

type OrderMine struct {
	OrderHST `bson:",inline"`
	TypeWork int    `bson:"type_work" json:"type_work"`
	Address  string `bson:"address" json:"address"`
}

func GetOrderMine(start int64, end int64, empID string) (ordMine []*OrderMine, err error) {
	fmt.Printf(" MINE START", start)
	fmt.Printf(" MINE END", end)
	fmt.Printf(" EMP", empID)
	var query []bson.M
	var queryMatch = bson.M{}
	var status = string(common.ORDER_STATUS_FINISHED)
	queryMatch["item_status"] = status
	queryMatch["emp_id"] = empID
	queryMatch["created_at"] = bson.M{"$lte": end, "$gte": start}
	var joinOrder = bson.M{
		"from":         "order",
		"localField":   "order_id",
		"foreignField": "_id",
		"as":           "order"}
	var sortDate = bson.M{"created_at": -1}
	var project = bson.M{
		"created_at":     1,
		"order_id":       1,
		"item_id":        1,
		"hour_day":       1,
		"money_day":      1,
		"item_status":    1,
		"start_work_day": 1,
		"end_work_day":   1,
		"address":        "$order.address_loc.address",
		"type_work":      "$order.type_work",
	}
	query = []bson.M{
		{"$match": queryMatch},
		{"$sort": sortDate},
		{"$lookup": joinOrder},
		{"$unwind": "$order"},
		{"$project": project},
	}
	err = OrderHstTable.Pipe(query).All(&ordMine)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	}
	return ordMine, err
}
func GetOrderHistory() ([]*OrderHST, error) {
	var ordersHistory []*OrderHST
	err := OrderHstTable.FindWhere(bson.M{}, &ordersHistory)
	return ordersHistory, err
}
