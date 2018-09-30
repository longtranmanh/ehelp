package order

import (
	"ehelp/common"
	//"ehelp/o/service"
	"ehelp/o/tool"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"sort"
	// "strings"
	// "ehelp/o/user/customer"
	// "ehelp/o/user/employee"
	//"time"
)

func GetAllOrderOpenAccept() ([]*Order, error) {
	var orders []*Order
	var timeNow = common.GetTimeNowVietNam().Unix()
	var err = OrderTable.FindWhere(bson.M{
		"status": bson.M{"$in": []string{string(common.ORDER_STATUS_OPEN), string(common.ORDER_STATUS_ACCEPTED), string(common.ORDER_STATUS_BIDDING)}},
		"date_end": bson.M{
			"$gte": timeNow,
			"$lte": timeNow,
		},
	}, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

type OrderTimePush struct {
	*Order
	TimePushStart float32
	HourStartItem float32
	IsUsed        bool
}

type OrderCheckChange struct {
	*Order
	HourStartItem float32
	HourEndItem   float32
	StatusItem    common.ItemOrderStatus
	IsUsedBf      bool
	IsUsedEnd     bool
	IsUsedMissed  bool
	Done          chan struct{}
}

func GetAllOrderAcceptWorking() (mapOrder map[int][]*OrderTimePush, err error) {
	var timeNow = common.GetTimeNowVietNam()
	var ords []*Order
	err = OrderTable.FindWhere(bson.M{
		"status": bson.M{"$in": []string{string(common.ORDER_STATUS_ACCEPTED), string(common.ORDER_STATUS_WORKING)}},
	}, &ords)
	if err != nil {
		return nil, err
	}
	mapOrder = make(map[int][]*OrderTimePush, 0)
	var orderAdd = make([]*OrderTimePush, 0)
	var weekDayNow = timeNow.AddDate(0, 0, 1).Weekday()
	for _, item := range ords {
		for _, itemWork := range item.DayWeeks {
			if common.ConvertTimeEpochToWeek(itemWork.DateIn) == weekDayNow {
				var hourItem = int(itemWork.HourStart)
				var ordTime = OrderTimePush{}
				ordTime.Order = item
				ordTime.TimePushStart = itemWork.HourStart
				if ords, ok := mapOrder[hourItem]; ok {
					ords = append(ords, &ordTime)
					mapOrder[hourItem] = ords
				} else {
					orderAdd = append(orderAdd, &ordTime)
					mapOrder[int(hourItem)] = orderAdd
				}
				break
			}
		}
	}
	return
}

func GetAllOrderCacheDay() (ords []*Order, err error) {
	var beginDay = common.BeginningOfDayVN().Unix()
	var endDay = common.EndOfDayVN().Unix()
	err = OrderTable.FindWhere(bson.M{
		"status": bson.M{"$in": []string{string(common.ORDER_STATUS_ACCEPTED), string(common.ORDER_STATUS_WORKING), string(common.ORDER_STATUS_BIDDING)}},
		"day_weeks.date_in": bson.M{
			"$gte": beginDay,
			"$lte": endDay,
		},
	}, &ords)
	return
}

func GetAllOrderBiddingDay() (ords []*Order, err error) {
	err = OrderTable.Find(bson.M{
		"status": common.ORDER_STATUS_BIDDING,
	}).Sort("day_start_work").All(&ords)
	return
}

func GetOrderById(idOrder string) (*Order, error) {
	var order *Order
	var err = OrderTable.FindByID(idOrder, &order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

type CustomerEmp struct {
	LinkAvatar string `bson:"link_avatar" json:"link_avatar"`
	FullName   string `bson:"full_name" json:"full_name"`
	Phone      string `bson:"phone" json:"phone"`
	Address    string `bson:"address" json:"address"`
	Email      string `bson:"email" json:"email"`
	ID         string `json:"id" bson:"_id"`
}

type OrderCusPromotion struct {
	Order `bson:",inline"`
	//Promotion *promotion.Promotion `json:"promotion" gorethink:"promotion"`
	ToolService []*tool.Tool `bson:"tools" json:"tools"`
}

type OrderAll struct {
	*Order      `bson:",inline"`
	ToolService []*tool.Tool `bson:"tools" json:"tools"`
}

func GetOrderAndCusAndEmp(idOrder string, typeGet int) (*OrderCusPromotion, error) {
	var query = []bson.M{}
	var queryMatch = bson.M{"_id": idOrder}

	// var joinCus = bson.M{
	// 	"from":         "customer",
	// 	"localField":   "cus_id",
	// 	"foreignField": "_id",
	// 	"as":           "customer",
	// }
	// var unWindCus = "$customer"

	// var joinEmp = bson.M{
	// 	"from":         "employee",
	// 	"localField":   "emp_id",
	// 	"foreignField": "_id",
	// 	"as":           "employee",
	// }
	// var unWindEmp = "$employee"

	switch typeGet {
	case 1:
		query = []bson.M{
			{"$match": queryMatch},
			//{"$lookup": joinCus},
			//{"$unwind": unWindCus},
		}
	case 2:
		query = []bson.M{
			{"$match": queryMatch},
			//	{"$lookup": joinEmp},
			//{"$unwind": unWindEmp},
		}
	case 3:
		query = []bson.M{
			{"$match": queryMatch},
			//	{"$lookup": joinCus},
			//	{"$unwind": unWindCus},
			//	{"$lookup": joinEmp},
			//	{"$unwind": unWindEmp},
		}
	}

	var ordPrm *OrderCusPromotion
	err := OrderTable.Pipe(query).One(&ordPrm)
	if err != nil {
		if err.Error() == common.NOT_EXIST {
			return nil, errors.New("Đơn hàng không tồn tại!")
		}
	}
	return ordPrm, nil
}

type Rest struct {
	ID string `bson:"_id" json:"_id"`
}

func GetAllOrderToExpired() ([]*Rest, error) {
	var allOrder []*Rest
	var query []bson.M
	var queryMatch = bson.M{}
	queryMatch["status"] = bson.M{"$in": []string{string(common.ORDER_STATUS_BIDDING), string(common.ORDER_STATUS_ACCEPTED)}}
	queryMatch["day_start_work"] = bson.M{"$lte": common.GetTimeNowVietNam().Unix()}
	query = []bson.M{
		{"$match": queryMatch},
		{"$project": bson.M{
			"_id": 1,
			"last_day_week": bson.M{
				"$arrayElemAt": []interface{}{"$day_weeks.date_in", -1},
			}},
		},
		{"$match": bson.M{"last_day_week": bson.M{
			"$lte": common.GetTimeNowVietNam().Unix(),
		}}},
	}
	err := OrderTable.Pipe(query).All(&allOrder)

	if err != nil && err.Error() == common.NOT_EXIST {
		return nil, nil
	}
	return allOrder, err
}

func GetAllOrderToFinished() ([]*Order, error) {
	var allOrder []*Order
	err := OrderTable.FindWhere(bson.M{
		"status": common.ORDER_STATUS_WORKING,
		"date_end": bson.M{
			"$lt": common.GetTimeNowVietNam().Unix(),
		},
	}, &allOrder)
	if err != nil && err.Error() == common.NOT_EXIST {
		return nil, nil
	}
	return allOrder, err
}

func (ord *Order) CheckAppendOrderEmp(empID string) error {
	fmt.Printf("VO CheckAppendOrderEmp", "")
	var status = []string{string(common.ORDER_STATUS_ACCEPTED), string(common.ORDER_STATUS_WORKING)}
	var ords, err = GetListOrderUserByStatus(empID, status, 1)
	if err != nil {
		return err
	}
	var hourNow = common.HourMinute()
	if ords != nil && len(ords) > 0 {
		sort.Sort(common.DayWeeks(ord.DayWeeks))
		if hourNow > ord.DayWeeks[0].HourStart && ord.TypeWork == common.TYPE_ONE_WEEK && common.CompareDayTime(common.GetTimeNowVietNam(), ord.DayWeeks[0].DateIn) == 0 {
			return errors.New("Đơn đã quá thời gian nhận việc!")
		}
		for _, dayNew := range ord.DayWeeks {
			for _, itemOld := range ords {
				sort.Sort(common.DayWeeks(itemOld.DayWeeks))
				for _, dayOld := range itemOld.DayWeeks {
					if common.CompareDayByTimeInt(dayNew.DateIn, dayOld.DateIn) == 0 {
						if (dayNew.HourEnd >= dayOld.HourStart && dayNew.HourEnd <= dayOld.HourEnd) ||
							(dayNew.HourStart >= dayOld.HourStart && dayNew.HourStart <= dayOld.HourEnd) ||
							(dayOld.HourEnd >= dayNew.HourStart && dayOld.HourEnd <= dayNew.HourEnd) ||
							(dayOld.HourStart >= dayNew.HourStart && dayOld.HourStart <= dayNew.HourEnd) {
							return errors.New("Đơn trùng lịch làm việc!")
						}
						if (dayOld.HourStart > dayNew.HourEnd && dayOld.HourStart-dayNew.HourEnd < 1) ||
							(dayOld.HourEnd > dayNew.HourStart && dayOld.HourEnd-dayNew.HourStart < 1) ||
							(dayNew.HourStart > dayOld.HourEnd && dayNew.HourStart-dayOld.HourEnd < 1) ||
							(dayNew.HourEnd > dayOld.HourStart && dayNew.HourEnd-dayOld.HourStart < 1) {
							return errors.New("Đơn phải cách tối thiểu 1h so với đơn cũ!")
						}
					}

				}
			}
		}
	}
	fmt.Printf("Khoong loi", "")
	return nil
}

func GetListOrderUserByStatus(userId string, status []string, role int) (ords []*Order, err error) {
	var queryMatch = bson.M{}
	queryMatch["status"] = bson.M{"$in": status}
	if role == 1 {
		queryMatch["emp_id"] = userId
	} else {
		queryMatch["cus_id"] = userId
	}
	err = OrderTable.FindWhere(queryMatch, &ords)
	if err != nil && err.Error() == common.NOT_EXIST {
		return nil, nil
	}
	return ords, err
}

func GetListOrderByStatus(userId string, serviceEmps []string, addressEmp string, role int, serviceId string, status []string, skip int, limit int) (ords []*OrderCusPromotion, err error) {
	var query []bson.M
	var queryMatch = bson.M{}
	var timeNow = common.GetTimeNowVietNam().Unix()
	queryMatch["status"] = bson.M{"$in": status}
	var statusBidding = string(common.ORDER_STATUS_BIDDING)

	var sortDate = bson.M{"day_start_work": -1}
	// var joinService = bson.M{
	// 	"from":         "service",
	// 	"localField":   "service_works",
	// 	"foreignField": "_id",
	// 	"as":           "services"}
	var joinTool = bson.M{
		"from":         "tool",
		"localField":   "tool_services",
		"foreignField": "_id",
		"as":           "tools"}

	if role == 2 { //Employee
		if status[0] != statusBidding {
			queryMatch["emp_id"] = userId

		} else {
			// if len(addressEmp) > 0 {
			// 	queryMatch["address_loc.address"] = bson.M{"$regex": addressEmp, "$options": "i"}
			// }
			if len(serviceEmps) > 0 {
				queryMatch["service_works"] = bson.M{"$in": serviceEmps}
				queryMatch["day_start_work"] = bson.M{"$gte": timeNow + 10800}
			}
		}
		// var joinCus = bson.M{
		// 	"from":         "customer",
		// 	"localField":   "cus_id",
		// 	"foreignField": "_id",
		// 	"as":           "customer",
		// }
		//var unWindCus = bson.M{"path": "$customer", "preserveNullAndEmptyArrays": true}
		query = []bson.M{
			{"$match": queryMatch},
			{"$sort": sortDate},
			{"$skip": skip},
			{"$limit": limit},
			// {"$lookup": joinCus},
			// {"$lookup": joinService},
			{"$lookup": joinTool},
			//{"$unwind": unWindCus},
		}
	} else {
		// var joinCus = bson.M{
		// 	"from":         "employee",
		// 	"localField":   "emp_id",
		// 	"foreignField": "_id",
		// 	"as":           "employee",
		// }
		// var unWindEmp = bson.M{"path": "$employee", "preserveNullAndEmptyArrays": true} // checknull
		queryMatch["cus_id"] = userId
		query = []bson.M{
			{"$match": queryMatch},
			{"$sort": sortDate},
			{"$skip": skip},
			{"$limit": limit},
			//{"$lookup": joinCus},
			//{"$lookup": joinService},
			{"$lookup": joinTool},
			//{"$unwind": unWindEmp},
		}
	}
	err = OrderTable.Pipe(query).All(&ords)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	}
	return ords, err
}

func CheckCustomerNewOfEmployee(cusID string, empID string) (countOrder int, err error) { // khi làm việc
	countOrder, err = OrderTable.CountWhere(bson.M{
		"cus_id": cusID,
		"emp_id": empID,
	})
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	}
	return
}
