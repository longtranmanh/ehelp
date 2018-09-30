package system

import (
	"ehelp/cache"
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/order_hst"
	"ehelp/o/push_token"
	oAuth "ehelp/o/user/auth"
	"ehelp/x/fcm"
	"fmt"
)

func canceled(ord *order.Order) {
	CreateOrderHst(ord.CusID, ord.ID, ord.ServiceWorks, common.ORDER_STATUS_CANCELED)
}
func accepted(ord *order.Order) {
	fmt.Printf("========= EmpID: "+ord.EmpID+" ID: "+ord.ID, ord.ServiceWorks)
	CreateOrderHst(ord.EmpID, ord.ID, ord.ServiceWorks, common.ORDER_STATUS_ACCEPTED)
	var emp, _ = cache.GetEmpID(ord.EmpID)
	fmt.Println("==== KHÁCH" + ord.CusID)
	pushs, _ := push_token.GetPushsUserId(ord.CusID) // phải distinic
	if len(pushs) > 0 {
		var noti = fcm.FmcMessage{
			Title: "Đã có người nhận!",
			Body:  emp.FullName + " đã nhận đơn"}
		fcm.FcmCustomer.SendToMany(pushs, noti)
		fmt.Println("====" + emp.FullName + " đã nhận đơn")
	}
	fmt.Printf("QUA ALL")
}
func bidding(ord *order.Order) {
	var empIds, _ = oAuth.GetListEmpVsOrderBidding(ord.ServiceWorks, ord.AddressLoc.Address)
	var pushs, err2 = push_token.GetPushsUserIds(empIds) // phải distinic
	logAction.Errorf("push_token.GetPushsUserIds", err2)
	var noti = fcm.FmcMessage{
		Title: "Có việc mới!",
		Body:  "Công việc tại " + ord.AddressLoc.Address}
	fcm.FcmEmployee.SendToMany(pushs, noti)
	fmt.Printf("========= CUSID: "+ord.CusID+" ID: "+ord.ID, ord.ServiceWorks)
	CreateOrderHst(ord.CusID, ord.ID, ord.ServiceWorks, common.ORDER_STATUS_BIDDING)
	fmt.Printf("QUA ALL")
}

func working(ord *order.Order, itemOrder *common.DayWeek) {
	CreateItemOrderHst(itemOrder.HourDay, float32(ord.PriceEnd),
		ord.CusID, ord.EmpID, ord.ServiceWorks, ord.ID, common.ORDER_STATUS_WORKING, itemOrder.IdItem,
		common.ITEM_ORDER_STATUS_WORKING, itemOrder.HourStart, itemOrder.HourEnd)
	var pushs, _ = push_token.GetPushsUserId(ord.CusID) // phải distinic
	var empOrd, _ = cache.GetEmpID(ord.EmpID)
	var noti = fcm.FmcMessage{
		Title: "Bắt đầu làm việc!",
		Body:  empOrd.FullName + " vừa bắt đầu làm việc!"}
	fcm.FcmCustomer.SendToMany(pushs, noti)
	fmt.Println("====" + empOrd.FullName + " đã nhận đơn")
}

func finished(ord *order.Order, itemOrder *common.DayWeek) {
	CreateItemOrderHst(itemOrder.HourDay, float32(ord.PriceEnd),
		ord.CusID, ord.EmpID, ord.ServiceWorks, ord.ID, common.ORDER_STATUS_WORKING, itemOrder.IdItem,
		common.ITEM_ORDER_STATUS_FINISHED, itemOrder.HourStart, itemOrder.HourEnd)
	var countCusNew, _ = order.CheckCustomerNewOfEmployee(ord.CusID, ord.EmpID)
	oAuth.UpdateCusNewAndHour(ord.EmpID, countCusNew, itemOrder.HourDay)
	var pushs, _ = push_token.GetPushsUserId(ord.CusID) // phải distinic
	var empOrd, _ = cache.GetEmpID(ord.EmpID)
	if ord.Status == common.ORDER_STATUS_FINISHED {
		var noti = fcm.FmcMessage{
			Title: "Công việc đã hoàn thành!",
			Body:  empOrd.FullName + " đã hoàn thành đầy đủ đơn mà bạn đặt! Lên đơn mới nếu muốn tìm người giúp việc!"}
		fcm.FcmCustomer.SendToMany(pushs, noti)
	} else {
		var noti = fcm.FmcMessage{
			Title: "Công việc đã hoàn thành!",
			Body:  empOrd.FullName + " đã hoàn thành việc ngày hôm nay!"}
		fcm.FcmCustomer.SendToMany(pushs, noti)
	}
}

func CreateOrderHst(userId string, orderID string, services []string, statusOrder common.OrderStatus) {
	var ordHst = order_hst.OrderHST{
		CusId:       userId,
		Services:    services,
		OrderId:     orderID,
		OrderStatus: statusOrder,
		ItemStatus:  common.ITEM_ORDER_STATUS_NEW,
	}
	ordHst.CrateOrderHistory()
}
func CreateItemOrderHst(itemHour float32, itemMoney float32, cusId string, empId string, services []string, orderID string, statusOrder common.OrderStatus, itemId string, statusItem common.ItemOrderStatus, startDay float32, endDay float32) {
	var ordHst = order_hst.OrderHST{
		CusId:        cusId,
		EmpId:        empId,
		Services:     services,
		ItemId:       itemId,
		ItemStatus:   statusItem,
		OrderId:      orderID,
		OrderStatus:  statusOrder,
		HourDay:      itemHour,
		MoneyDay:     itemMoney,
		StartWorkDay: startDay,
		EndWorkDay:   endDay,
	}
	ordHst.CrateOrderHistory()
}
