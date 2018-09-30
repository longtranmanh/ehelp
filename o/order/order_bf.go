package order

import (
	"fmt"
	"sort"
	"time"
	//"ehelp/o/user"
	//"ehelp/x/rest"
	"ehelp/o/service"
	"errors"
	validator "gopkg.in/go-playground/validator.v9"
	// "gopkg.in/mgo.v2/bson"
	//"strings"
	"ehelp/common"
	"ehelp/x/rest"
)

var validate = validator.New()

func (ord *Order) create() error {
	//check id service => kèm theo service //check id tool
	var sers, errs = service.GetServiceAndTool(ord.ServiceWorks)
	if errs != nil {
		return errs
	}

	ord.Status = common.ORDER_STATUS_BIDDING
	// tính tiền giờ
	var hourAll, priceAllHour, priceTool, priceEnd, err = ord.MathPriceOrder.MathPriceOrder()
	if err != nil {
		return err
	}
	var serDs = make([]*ServiceDetail, 0)
	for _, ser := range sers {
		var serD = &ServiceDetail{
			ID:           ser.ID,
			Name:         ser.Name,
			NodeServices: ser.NodeServices,
		}
		serDs = append(serDs, serD)
	}
	ord.Services = serDs
	ord.AllHourWork = hourAll
	ord.PriceAllHour = priceAllHour
	ord.PriceTool = priceTool
	ord.PriceEnd = int(priceEnd)
	if err := validate.Struct(ord); err != nil {
		return rest.BadRequestValid(err)
	}
	return nil
}

func (ord *Order) update(status common.OrderStatus) error {
	return ord.CheckStatus(status, "", "")
}

func (ord *Order) CheckStatus(status common.OrderStatus, empId string, cusID string) (err error) {
	var statusOrder = ord.Status
	switch statusOrder {
	case common.ORDER_STATUS_BIDDING:
		if status == common.ORDER_STATUS_FINISHED || status == common.ORDER_STATUS_WORKING {
			err = errors.New("Đang tìm người làm!")
		}
	case common.ORDER_STATUS_ACCEPTED:
		if len(cusID) > 0 {
			if status == common.ORDER_STATUS_FINISHED {
				err = errors.New("Bạn chỉ có thể hủy đơn!")
			}
		} else if status == common.ORDER_STATUS_FINISHED {
			err = errors.New("Bạn chưa làm việc trước khi kết thúc!")
		} else if status == common.ORDER_STATUS_ACCEPTED {
			err = errors.New("Đã có người nhận đơn!")
		}
	case common.ORDER_STATUS_OPEN:
		if status != common.ORDER_STATUS_BIDDING {
			err = errors.New("Đơn hàng đã hết hạn!")
		}
	case common.ORDER_STATUS_WORKING:
		if status == common.ORDER_STATUS_CANCELED || status == common.ORDER_STATUS_BIDDING || status == common.ORDER_STATUS_OPEN {
			err = errors.New("Đã có người làm việc!")
		}
		if status == common.ORDER_STATUS_FINISHED && len(cusID) > 0 && ord.CheckItemWorking() {
			err = errors.New("Đang có người làm việc! Hãy kết thúc khi nhân viên hoàn thành!")
		}
	case common.ORDER_STATUS_FINISHED:
		err = errors.New("Đơn đã kết thúc!")
	case common.ORDER_STATUS_CANCELED:
		if status != common.ORDER_STATUS_FINISHED && len(cusID) > 0 {
			err = errors.New("Không thể kết thúc khi chưa làm việc!")
		}
	}
	if statusOrder != common.ORDER_STATUS_BIDDING && statusOrder != common.ORDER_STATUS_CANCELED && empId != "" && ord.EmpID != empId {
		err = rest.Unauthorized("Đơn này không phải của bạn!")
		return
	}

	return
}

func (ord *Order) CheckTimeUpdateItem(isWorking bool) (itemCheck *common.DayWeek, err error) {
	var timeNow = common.GetTimeNowVietNam()
	fmt.Printf("== HE THONG :", timeNow)
	for _, item := range ord.MathPriceOrder.DayWeeks {
		if common.CompareDayTime(timeNow, item.DateIn) == 0 {
			if isWorking {
				var timeNowHour = common.HourMinuteEpoch(timeNow.Unix())
				var res = item.HourStart - timeNowHour
				fmt.Printf("== CHECK 30p :", res)
				if (res <= 0.5 && res >= 0) || (res <= 0 && res >= -0.5) {
					itemCheck = item
					break
				} else {
					var msgErr = "Giờ làm của đơn: " + common.ConvertF32ToString(item.HourStart)
					msgErr += "\nThời gian bắt đầu phải sớm hoặc muộn hơn trong khoảng 30 phút!"
					err = rest.WrapBadRequest(errors.New(msgErr), "")
					return
				}
			} else {
				itemCheck = item
				break
			}
		} else if item.Status == common.ITEM_ORDER_STATUS_WORKING && !isWorking {
			itemCheck = item
			break
		}
	}
	if itemCheck == nil {
		err = rest.WrapBadRequest(errors.New("Không có lịch làm việc hôm nay!"), "")
	}
	return
}

func (ord *Order) CheckItemWorking() (isWorking bool) {

	var dayNow = common.GetTimeNowVietNam().Day()
	for _, item := range ord.DayWeeks {
		var dateItem = time.Unix(item.DateIn, 0).Day()
		if item.Status == common.ITEM_ORDER_STATUS_WORKING && dayNow == dateItem {
			isWorking = true
			break
		}
	}
	return
}

func (ord *Order) CheckItemFinished() (isFinished bool) {
	sort.Sort(ord.DayWeeks)
	var item = ord.DayWeeks[len(ord.DayWeeks)-1]
	var dateItem = time.Unix(item.DateIn, 0).Day()
	var dayNow = common.GetTimeNowVietNam().Day()
	if dayNow == dateItem {
		isFinished = true
	}
	return
}
