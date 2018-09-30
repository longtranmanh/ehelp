package system

import (
	"ehelp/cache"
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/push_token"
	"ehelp/x/fcm"
	"ehelp/x/mlog"
	"ehelp/x/rest"
	"fmt"
)

var logSys = mlog.NewTagLog("System")

/*Get Order ID*/
func GetOrderID(id string) (*order.Order, bool) {
	if val, ok := CacheOrderByDay.Orders[id]; ok {
		return val.Order, true
	}
	var ord, err = order.GetOrderById(id)
	rest.AssertNil(err)
	// if ord != nil {
	// 	var weekDayNow = common.GetTimeNowVietNam().Weekday()
	// 	for _, itemWork := range ord.DayWeeks {
	// 		if common.ConvertTimeEpochToWeek(itemWork.DateIn) == weekDayNow {
	// 			var or = order.OrderCheckChange{
	// 				Order:         ord,
	// 				HourStartItem: itemWork.HourStart,
	// 				HourEndItem:   itemWork.HourEnd,
	// 				StatusItem:    itemWork.Status,
	// 			}
	// 			CacheOrderByDay.Orders[id] = &or
	// 			break
	// 		}
	// 	}
	// }
	return ord, false
}
func Launch() {
	cache.SetCacheCus()
	cache.SetCacheEmp()
	SetCacheOrderDay()
}

func GetSendPushWork() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var res = int(ord.HourStartItem) - common.GetTimeNowVietNam().Hour()
		if !ord.IsUsedBf && 1 == res && (ord.Status == common.ORDER_STATUS_ACCEPTED || ord.Status == common.ORDER_STATUS_WORKING) {
			ords = append(ords, ord)
		}
	}
	return
}

func GetSendPushWorkEnd() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var val = ord.HourEndItem - common.HourMinute()
		if !ord.IsUsedEnd && ord.StatusItem == common.ITEM_ORDER_STATUS_WORKING && val >= 0 && val <= 0.25 {
			ords = append(ords, ord)
		}
	}
	return
}

func checkAndSendPushBf() {
	var ordPushWork = GetSendPushWork()
	fmt.Printf("checkAndSendPushBf", len(ordPushWork))
	if len(ordPushWork) > 0 {
		for _, ord := range ordPushWork {
			var body = "Thời gian: "
			var notify = fcm.FmcMessage{
				Title: "Hôm nay có lịch làm việc!",
				Body: body + common.ConvertF32ToString(ord.HourStartItem) +
					".\nĐịa chỉ: " + ord.AddressLoc.Address +
					".\nVui lòng đến đúng giờ!",
			}
			var pushs, _ = push_token.GetPushsUserId(ord.EmpID)
			fcm.FcmEmployee.SendToMany(pushs, notify)
			ord.IsUsedBf = true
		}

	}
}

func checkAndSendPushEnd() {
	var ordPushWork = GetSendPushWorkEnd()
	if len(ordPushWork) > 0 {
		for _, ord := range ordPushWork {
			var body = "Thời gian: "
			var notify = fcm.FmcMessage{
				Title: "Sắp hết giờ!",
				Body:  body + common.ConvertF32ToString(ord.HourEndItem) + " sẽ đủ số giờ làm. Hãy bấm kết thúc!",
			}
			var pushs, _ = push_token.GetPushsUserId(ord.EmpID)
			fcm.FcmEmployee.SendToMany(pushs, notify)
			ord.IsUsedEnd = true
		}
	}
}

func getAutoMissed() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var val = ord.HourStartItem - common.HourMinute()
		if !ord.IsUsedMissed && ord.Status == common.ORDER_STATUS_BIDDING && (val >= 0 && val <= 3) {
			ords = append(ords, ord)
		}
	}
	fmt.Println("VO AUTO MISSED", len(ords))
	return
}

func SendPushAndChangeMissed() {
	var ordMisseds = getAutoMissed()
	var ordIds = make([]string, len(ordMisseds))
	for i, ord := range ordMisseds {
		fmt.Printf("ĐƠN HỦY: "+ord.ID, "Nguoi tao: "+ord.CusID)
		ordIds[i] = ord.ID
		var body = "Đơn tại: "
		var notify = fcm.FmcMessage{
			Title: "Đơn được hủy!",
			Body: body + ord.AddressLoc.Address + ".\nThời gian " +
				common.ConvertF32ToString(ord.HourStartItem) +
				" sẽ được hủy do không tìm được người làm.\nQuý khách vui lòng lên lại đơn để tìm người giúp việc!",
		}
		var pushs, _ = push_token.GetPushsUserId(ord.CusID)
		fcm.FcmCustomer.SendToMany(pushs, notify)
		ord.IsUsedMissed = true
		//delete(CacheOrderByDay.Orders, ord.ID)
	}
	order.UpdateStatusByIds(ordIds, common.ORDER_STATUS_OPEN)
}
