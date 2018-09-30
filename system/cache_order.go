package system

import (
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/x/mrw/event"
	"fmt"
	"time"
)

type CacheOrderWorker struct {
	Orders      map[string]*order.OrderCheckChange
	OrderUpdate chan *OrderAction
	doneAction  *event.Hub
}

func (tw *CacheOrderWorker) TriggerTicketAction(action *OrderAction) {
	tw.OrderUpdate <- action
}

var CacheOrderByDay = NewCacheOrderWorker()

func NewCacheOrderWorker() *CacheOrderWorker {
	return &CacheOrderWorker{
		Orders:      make(map[string]*order.OrderCheckChange, 0),
		OrderUpdate: make(chan *OrderAction, event.MediumHub),
	}
}

func (tw *CacheOrderWorker) OnActionDone() (event.Line, event.Cancel) {
	return tw.doneAction.NewLine()
}

func SetCacheOrderDay() {
	if CacheOrderByDay != nil && CacheOrderByDay.Orders != nil && len(CacheOrderByDay.Orders) > 0 {
		for k, _ := range CacheOrderByDay.Orders {
			delete(CacheOrderByDay.Orders, k)
		}
	}
	var timeNow = time.Now()
	var orders, _ = order.GetAllOrderCacheDay()
	fmt.Printf("SỐ ORDER", len(orders))
	if len(orders) > 0 {
		for _, item := range orders {
			for _, itemWork := range item.DayWeeks {
				if common.CompareDayTime(timeNow, itemWork.DateIn) == 0 {
					var or = order.OrderCheckChange{
						Order:         item,
						HourStartItem: itemWork.HourStart,
						HourEndItem:   itemWork.HourEnd,
						StatusItem:    itemWork.Status,
					}
					CacheOrderByDay.Orders[item.ID] = &or
					break
				}
			}
		}
	}
	fmt.Printf("Cache Day", len(CacheOrderByDay.Orders))
}
