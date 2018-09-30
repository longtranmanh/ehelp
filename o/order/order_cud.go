package order

import (
	"ehelp/common"
	"ehelp/x/mlog"
	"ehelp/x/rest"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

var logOrder = mlog.NewTagLog("Order")

func (ord *Order) CrateOrder() (*Order, error) {
	var err = ord.create()
	if err != nil {
		logOrder.Errorf("CrateOrder", err)
		return nil, err
	}
	err = OrderTable.Create(ord)
	return ord, err
}

func (ord *Order) UpdateOrder(status common.OrderStatus) (*Order, error) {
	ord.update(status)
	var err = OrderTable.UpdateId(ord.ID, ord)
	return ord, err
}

func UpdateByRate(orderID string) error {
	var newUp = map[string]interface{}{
		"rated": true,
	}
	return OrderTable.UpdateId(orderID, bson.M{
		"$set": newUp,
	})
}

func UpdateStatusOrders(ordIDs []*Rest, status string) (error, int) {
	var ids = make([]string, len(ordIDs))
	for i, item := range ordIDs {
		ids[i] = item.ID
	}
	var newUp = map[string]interface{}{
		"status": status,
	}
	var rest, err = OrderTable.UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": newUp})
	return err, rest.Updated
}

func UpdateStatusByIds(ordIDs []string, status common.OrderStatus) (error, int) {
	var newUp = map[string]interface{}{
		"status": status,
	}
	var rest, err = OrderTable.UpdateAll(bson.M{"_id": bson.M{"$in": ordIDs}}, bson.M{"$set": newUp})
	return err, rest.Updated
}

func (ord *Order) UpdateStatusOrder(status common.OrderStatus, empId string, cusID string, empCus *CustomerEmp) (*Order, error) {
	fmt.Print("VO UpdateStatusOrder")
	var err = ord.CheckStatus(status, empId, cusID)
	if err != nil {
		err = rest.BadRequestValid(err)
		return nil, err
	}
	fmt.Print("QUA UpdateStatusOrder")
	var newUp = map[string]interface{}{
		"status": status,
	}
	if len(empId) > 0 {
		newUp["emp_id"] = empId
	}
	if empCus != nil {
		newUp["employee"] = empCus
	}
	var tmNow = common.GetTimeNowVietNam().Unix()
	newUp["updated_at"] = tmNow
	err = OrderTable.UpdateId(ord.ID, bson.M{
		"$set": newUp,
	})
	if err == nil {
		if len(empId) > 0 {
			ord.EmpID = empId
		}
		if empCus != nil {
			ord.Employee = empCus
		}
		ord.Status = status
		ord.UpdatedAt = tmNow
	} else {
		logOrder.Errorf("UpdateStatusOrder", err)
	}
	return ord, err
}

func (ord *Order) UpdateStatusItem(itemID string, status common.ItemOrderStatus) (ordResp *Order, err error) {
	var items = map[string]interface{}{
		"day_weeks": ord.DayWeeks,
	}
	var statusOrder common.OrderStatus
	switch ord.Status {
	case common.ORDER_STATUS_ACCEPTED:
		if status == common.ITEM_ORDER_STATUS_WORKING {
			items["status"] = common.ORDER_STATUS_WORKING
			statusOrder = common.ORDER_STATUS_WORKING
		}
	case common.ORDER_STATUS_WORKING:
		if status == common.ITEM_ORDER_STATUS_FINISHED {
			items["status"] = common.ORDER_STATUS_FINISHED
		}

		if status == common.ITEM_ORDER_STATUS_FINISHED && ord.TypeWork == common.TYPE_ONE_WEEK && ord.CheckItemFinished() {
			statusOrder = common.ORDER_STATUS_FINISHED
		}
	case common.ORDER_STATUS_FINISHED:
		err = rest.WrapBadRequest(errors.New("Công việc đã kết thúc"), "")
		return
	case common.ORDER_STATUS_CANCELED:
		err = rest.WrapBadRequest(errors.New("Công việc đã bị hủy"), "")
		return
	}
	var timeNow = common.GetTimeNowVietNam()
	for _, item := range ord.DayWeeks {
		if itemID == item.IdItem {
			var timeWork = int64((item.HourEnd - item.HourStart) * 3600)
			var timeWork1 = timeNow.Unix() - item.MTime
			if status == common.ITEM_ORDER_STATUS_FINISHED && timeWork-timeWork1 > 0 {
				var mes = "Thời gian làm việc phải đủ " + common.ConvertF32ToString(item.HourDay)
				err = rest.WrapBadRequest(errors.New(mes), "")
				fmt.Printf("LỖI KẾT THÚC", err)
				return
			}
			if item.Status == status && common.CompareDayTime(timeNow, item.MTime) == 0 {
				err = rest.WrapBadRequest(errors.New("Bạn đã thực hiện thao tác này"), "")
				return
			}
			item.Status = status
			item.MTime = timeNow.Unix()
			break
		}
	}
	fmt.Println("")
	fmt.Println("QUA HẾT LỖI")
	var iStatus = len(statusOrder)
	if iStatus > 0 {
		items["status"] = statusOrder
	}
	err = OrderTable.UpdateId(ord.ID, bson.M{
		"$set": items,
	})
	if err != nil {
		return
	}
	if iStatus > 0 {
		ord.Status = statusOrder
	}
	ordResp = ord
	return
}
