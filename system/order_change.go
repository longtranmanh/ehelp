package system

import (
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/x/mlog"
	"encoding/json"
	"errors"
	"fmt"
)

var logAction = mlog.NewTagLog("ControllerAction")

func (action *OrderAction) setOrder(isOld bool, ord *order.Order) {
	if isOld {
		var ordChange = order.OrderCheckChange{
			Order: ord,
		}
		action.Order = &ordChange
	}
}

func (action *OrderAction) HandlerAction() {
	var ord *order.Order
	var isOldCache bool
	if len(action.OrderID) > 0 {
		ord, isOldCache = GetOrderID(action.OrderID)
	}
	switch action.Action {
	case common.ORDER_STATUS_BIDDING:
		ord = action.Order.Order
		var _, err = ord.CrateOrder()
		if err != nil {
			action.SetError(err)
			return
		}
		bidding(ord)
	case common.ORDER_STATUS_ACCEPTED:
		var empID = action.EmpId
		fmt.Printf("HAND EMP ID", empID)
		err := ord.CheckAppendOrderEmp(empID)
		if err != nil {
			fmt.Printf("ord.CheckAppendOrderEmp(empID)", err)
			action.SetError(err)
			return
		}
		var extra = order.CustomerEmp{}
		err = json.Unmarshal(action.Extra, &extra)
		if err != nil {
			action.SetError(err)
			return
		}
		_, err1 := ord.UpdateStatusOrder(common.ORDER_STATUS_ACCEPTED, empID, "", &extra)
		if err1 != nil {
			fmt.Printf("ord.UpdateStatusOrder(common.ORDER_STATUS_ACCEPTED, empID, )", err)
			action.SetError(err1)
			return
		}
		UpdateStatusCache(action.OrderID, common.ITEM_ORDER_STATUS_NEW, common.ORDER_STATUS_ACCEPTED)
		accepted(ord)
		action.setOrder(isOldCache, ord)
	case common.ORDER_STATUS_OPEN:
	case common.ORDER_STATUS_WORKING:
		var itemOrder, err = ord.CheckTimeUpdateItem(true)
		if err != nil {
			action.SetError(err)
			return
		}
		_, err1 := ord.UpdateStatusItem(itemOrder.IdItem, common.ITEM_ORDER_STATUS_WORKING)
		if err1 != nil {
			action.SetError(err1)
			return
		}
		UpdateStatusCache(action.OrderID, common.ITEM_ORDER_STATUS_WORKING, common.ORDER_STATUS_WORKING)
		working(ord, itemOrder)
	case common.ORDER_STATUS_CANCELED:
		_, err := ord.UpdateStatusOrder(common.ORDER_STATUS_CANCELED, "", ord.CusID, nil)
		if err != nil {
			action.SetError(err)
			return
		}
		UpdateStatusCache(action.OrderID, common.ITEM_ORDER_STATUS_FINISHED, common.ORDER_STATUS_CANCELED)
		canceled(ord)
	case common.ORDER_STATUS_FINISHED:
		var itemOrder, err = ord.CheckTimeUpdateItem(false)
		if err != nil {
			action.SetError(err)
			return
		}
		_, err = ord.UpdateStatusItem(itemOrder.IdItem, common.ITEM_ORDER_STATUS_FINISHED)
		if err != nil {
			action.SetError(err)
			return
		}
		UpdateStatusCache(action.OrderID, common.ITEM_ORDER_STATUS_FINISHED, common.ORDER_STATUS_FINISHED)
		finished(ord, itemOrder)
	default:
		err := errors.New("No Action")
		action.SetError(err)
		logAction.Errorf("HandlerAction", err)
	}
}

func UpdateStatusCache(orderID string, status common.ItemOrderStatus, statusOrder common.OrderStatus) {
	if CacheOrderByDay.Orders != nil {
		if val, ok := CacheOrderByDay.Orders[orderID]; ok {
			val.StatusItem = status
			val.Status = statusOrder
			val.IsUsedEnd = true
		}
	}
}
