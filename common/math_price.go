package common

import (
	"ehelp/o/service"
	"ehelp/o/tool"
	"ehelp/x/rest"
	"ehelp/x/rest/math"
	"errors"
	"fmt"
	"sort"
)

type MathPriceOrder struct {
	TypeWork     TypeWork     `bson:"type_work" json:"type_work" validate:"required"`
	Promotions   []string     `bson:"promotions" json:"promotions"`
	ServiceWorks []string     `bson:"service_works" json:"service_works" validate:"required"`
	ToolServices []string     `bson:"tool_services" json:"tool_services"`
	DayWeeks     DayWeeks     `bson:"day_weeks" json:"day_weeks"`
	PeopleOther  *PeopleOther `bson:"people_other" json:"people_other"`
	DayStartWork int64        `bson:"day_start_work" json:"day_start_work" validate:"required"`
}

type PeopleOther struct {
	Phone string `bson:"phone" json:"phone"`
	Name  string `bson:"name" json:"name"`
}

type DayWeek struct {
	IdItem    string          `bson:"id_item" json:"id_item" validate:"required"`
	DateIn    int64           `bson:"date_in" json:"date_in" validate:"required"` // 2,,3,4,5,6,7,8
	HourStart float32         `bson:"hour_start" json:"hour_start" validate:"required"`
	HourEnd   float32         `bson:"hour_end" json:"hour_end" validate:"required"`
	HourDay   float32         `bson:"hour_day" json:"hour_day" validate:"required"`
	MTime     int64           `bson:"mtime" json:"mtime"`
	Status    ItemOrderStatus `bson:"status" json:"status" validate:"required"`
}

type DayWeeks []*DayWeek

func (a DayWeeks) Len() int           { return len(a) }
func (a DayWeeks) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DayWeeks) Less(i, j int) bool { return a[i].DateIn < a[j].DateIn }

func (ord *MathPriceOrder) MathHourWork() (float32, error) {
	var err error
	var timeNow = BeginningOfDayVN().Unix()
	var hourMinuteNow = HourMinute()
	// if ord.TypeWork == TYPE_ONE_WEEK {
	// 	if timeNow > timeOrdStart {
	// 		rest.AssertNil(rest.WrapBadRequest(errors.New("Không thể đặt trước ngày hôm nay!"), ""))
	// 	} else if timeNow == timeOrdStart {
	// 		sort.Sort(DayWeeks(ord.DayWeeks))
	// 		var timeStartMin = ord.DayWeeks[0].HourStart
	// 		var dateIn = ord.DayWeeks[0].DateIn
	// 		var timeItemMin = BeginningOfDayInt64(dateIn).Unix()
	// 		if timeItemMin == timeOrdStart && HourMinute() > timeStartMin {
	// 			rest.AssertNil(rest.WrapBadRequest(errors.New("Đặt lại! Giờ làm việc lớn hơn giờ hiện tại!"), ""))
	// 		}
	// 	}
	// }
	var hourInWeek float32
	var dayWeeks = DayWeeks{}
	sort.Sort(DayWeeks(ord.DayWeeks))

	for _, item := range ord.DayWeeks {
		var timeOrdStart = BeginningOfDayInt64VN(item.DateIn).Unix()
		if timeOrdStart < timeNow || (timeOrdStart == timeNow && hourMinuteNow > item.HourStart) {
			continue
		}

		item.IdItem = math.RandStringUpper("", 6)
		var hourDay = item.HourEnd - item.HourStart
		if hourDay < 2 {
			err = rest.BadRequestValid(errors.New("Số giờ tối thiểu 1 ngày là 2h!"))
			return 0, err
		}
		item.HourDay = hourDay
		item.Status = ITEM_ORDER_STATUS_NEW
		hourInWeek += hourDay
		dayWeeks = append(dayWeeks, item)
	}
	if len(dayWeeks) == 0 {
		err = rest.BadRequestValid(errors.New("Không được lên đơn trong quá khứ!"))
		return 0, err
	}
	ord.DayWeeks = dayWeeks
	fmt.Println(hourInWeek)
	return hourInWeek, err
}

func (ord *MathPriceOrder) MathPriceOrder() (hourAll float32, priceAllHour float32, priceTool float32, priceEnd float32, err error) {
	hourAll, err = ord.MathHourWork()
	if err != nil {
		return
	}
	switch ord.TypeWork {
	case TYPE_ONE_WEEK:
		//hourAll = hourAll * 1
	case TYPE_TWO_WEEK:
		//hourAll = hourAll * 2
	case TYPE_THREE_WEEK:
		//hourAll = hourAll * 3
	case TYPE_FOUR_WEEK:
	//	hourAll = hourAll * 4
	default:
		rest.AssertNil(rest.BadRequestValid(errors.New("Không tồn tại loại làm việc này!")))
	}
	if len(ord.ServiceWorks) > 0 {
		var service, err = service.GetByID(ord.ServiceWorks[0])
		fmt.Printf("=====Service =====", service)
		rest.AssertNil(err)
		priceAllHour = hourAll * float32(service.PricePerHour)
		fmt.Printf("priceAllHour", priceAllHour)
	}
	if len(ord.ToolServices) > 0 {
		var srcTools, err = tool.GetToolByArrayID(ord.ToolServices)
		if err != nil && err.Error() != NOT_EXIST {
			rest.AssertNil(rest.WrapBadRequest(err, ""))
		}
		for _, item := range srcTools {
			priceTool += float32(item.Price)
		}
	}
	priceEnd = priceAllHour + priceTool
	return
}
