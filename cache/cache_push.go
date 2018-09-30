package cache

// import (
// 	"ehelp/common"
// 	"ehelp/o/order"
// 	"ehelp/x/mrw/event"
// 	"ehelp/x/rest"
// )

// type OrderCache struct {
// 	*order.Order
// 	TimeStartDay float32
// 	TimeEndDay   float32
// }
// type CacheOrderWorker struct {
// 	cache       map[string]*OrderCache
// 	OrderUpdate chan *order.OrderCheckChange
// 	doneAction  *event.Hub
// }

// func NewCacheOrderWorker() *CacheOrderWorker {
// 	return &CacheOrderWorker{
// 		OrderUpdate: make(chan *order.OrderCheckChange, event.MediumHub),
// 	}
// }

// type CacheOrderWorker struct {
// 	Orders                map[string]*order.Order
// 	OrderBiddings         map[int][]*order.OrderCheckChange
// 	OrderWorkingAccepteds map[int][]*order.OrderCheckChange
// 	OrderUpdate           chan *order.OrderCheckChange
// 	doneAction            *event.Hub
// }

// var cacheOrderByDay = NewCacheOrderWorker()

// func NewCacheOrderWorker() *CacheOrderWorker {
// 	return &CacheOrderWorker{
// 		OrderUpdate:           make(chan *order.OrderCheckChange, event.MediumHub),
// 		OrderBiddings:         make(map[int][]*order.OrderCheckChange, 0),
// 		OrderWorkingAccepteds: make(map[int][]*order.OrderCheckChange, 0),
// 	}
// }

// func (tw *CacheOrderWorker) OnActionDone() (event.Line, event.Cancel) {
// 	return tw.doneAction.NewLine()
// }

// var CacheOrderBidding = make(map[string]*OrderCache)

// func GetOrderID(id string) *order.Order {
// 	if val, ok := CacheOrderBidding[id]; ok {
// 		return val.Order
// 	}
// 	var ord, err = order.GetOrderById(id)
// 	rest.AssertNil(err)
// 	var weekDayNow = common.GetTimeNowVietNam().Weekday()
// 	var start float32
// 	var end float32
// 	for _, itemWork := range ord.DayWeeks {
// 		if common.ConvertTimeEpochToWeek(itemWork.DateIn) == weekDayNow {
// 			start = itemWork.HourStart
// 			end = itemWork.HourEnd
// 			break
// 		}
// 	}
// 	var or = OrderCache{
// 		Order:        ord,
// 		TimeStartDay: start,
// 		TimeEndDay:   end,
// 	}
// 	CacheOrderBidding[ord.ID] = &or
// 	return ord
// }

// func RefeshBidding(orderId string, status common.OrderStatus) {
// 	switch status {
// 	case common.ORDER_STATUS_BIDDING:
// 	default:
// 		delete(CacheOrderBidding, orderId)
// 	}
// }

// func AddCacheListOrderDay() (ordCaches map[string]*OrderBidding) {
// 	var ords, _ = order.GetAllOrderBiddingDay()
// 	var lenOrds = len(ords)
// 	if ords != nil && lenOrds > 0 {
// 		ordCaches = make(map[string]*OrderBidding, lenOrds)
// 		var weekDayNow = common.GetTimeNowVietNam().AddDate(0, 0, 1).Weekday()
// 		for i := 0; i < lenOrds; i++ {
// 			for _, itemWork := range ords[i].DayWeeks {
// 				if common.ConvertTimeEpochToWeek(itemWork.DateIn) == weekDayNow {
// 					var ordBid = OrderBidding{
// 						Order:        ords[i],
// 						TimeEndDay:   itemWork.HourEnd,
// 						TimeStartDay: itemWork.HourStart,
// 					}
// 					ordCaches[ords[i].ID] = &ordBid
// 				}
// 			}
// 		}
// 	}
// 	return
// }
