package order_server

import (
	"ehelp/cache"
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/rate"
	oAuth "ehelp/o/user/auth"
	"ehelp/o/user/customer"
	"ehelp/o/user/employee"
	"ehelp/system"
	"ehelp/x/mrw/encode"
	"ehelp/x/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type OrderServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewOrderServer(parent *gin.RouterGroup, name string) {
	var s = &OrderServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("create", s.handleCreate)
	s.POST("status/accepted", s.handleAccepted)
	s.POST("status/working", s.handleWorking)
	//s.POST("status/cus_finished", s.handleCusFinished)
	s.POST("status/emp_finished", s.handleEmpFinished)
	s.POST("status/canceled", s.handleCancelled)
	s.POST("rate", s.handleRating)
	s.GET("list_order", s.handleListOrder)
	s.GET("order_mine", s.handleOrderMine)
}

func (s *OrderServer) handleCreate(ctx *gin.Context) {
	var cus = oAuth.GetCusFromToken(ctx.Request)
	var ord *order.Order
	rest.AssertNil(ctx.ShouldBindJSON(&ord))
	var cust = &order.CustomerEmp{
		ID:       cus.ID,
		Phone:    cus.Phone,
		FullName: cus.FullName,
	}
	ord.Customer = cust
	//var hourMoney = cache.GetCacheHourMoney()
	ord.CusID = cus.ID
	var action = system.NewOrderAction()
	action.Action = common.ORDER_STATUS_BIDDING
	var ordChange = order.OrderCheckChange{
		Order:      ord,
		StatusItem: common.ITEM_ORDER_STATUS_NEW,
	}
	action.Order = &ordChange
	system.CacheOrderByDay.TriggerTicketAction(action)
	var ordRes, err = action.Wait()
	rest.AssertNil(err)
	s.SendData(ctx, dataRespone(ordRes, nil, cus))
}

func dataRespone(ordRes *order.Order, emp *employee.Employee, cus *customer.Customer) *order.OrderAll {
	var ordCus = order.OrderAll{
		Order: ordRes,
	}
	return &ordCus
}

func (s *OrderServer) handleAccepted(ctx *gin.Context) {
	var emp = oAuth.GetEmpFromToken(ctx.Request)
	var body = struct {
		OrderID string `json:"order_id"`
	}{}
	ctx.BindJSON(&body)
	var ex = &order.CustomerEmp{
		FullName:   emp.FullName,
		LinkAvatar: emp.LinkAvatar,
		Phone:      emp.Phone,
		ID:         emp.ID,
		Email:      emp.Email,
		Address:    emp.Address,
	}
	var extra, _ = json.Marshal(ex)
	//ord := system.CacheOrderByDay.Orders[body.OrderID]
	ordRes := ActionChange(body.OrderID, emp.ID, common.ORDER_STATUS_ACCEPTED, extra)
	fmt.Print("QUA ERRR: " + ordRes.CusID)
	var cus, _ = cache.GetCusID(ordRes.CusID)
	fmt.Print("QUA CUS: " + ordRes.CusID)
	s.SendData(ctx, dataRespone(ordRes, emp, cus))
}

// func (s *OrderServer) handleCusFinished(ctx *gin.Context) {
// 	var body = struct {
// 		OrderID string `json:"order_id"`
// 	}{}
// 	ctx.BindJSON(&body)
// 	var ord = cache.MustGetByID(body.OrderID)
// 	var cus = oAuth.GetCusFromToken(ctx.Request)
// 	if cus.ID != ord.CusID {
// 		rest.AssertNil(rest.Unauthorized("Đơn không phải của bạn!"))
// 	}
// 	ord.UpdateStatusOrder(common.ORDER_STATUS_FINISHED, "", cus.ID)
// 	cache.Refresh(ord)
// 	var ordCus, err = order.GetOrderAndCusAndEmp(body.OrderID, 3)
// 	rest.AssertNil(err)
// 	s.SendData(ctx, ordCus)
// }

func (s *OrderServer) handleEmpFinished(ctx *gin.Context) {
	var body = struct {
		OrderID string `json:"order_id"`
	}{}
	ctx.BindJSON(&body)
	var emp = oAuth.GetEmpFromToken(ctx.Request)
	ordRes := ActionChange(body.OrderID, emp.ID, common.ORDER_STATUS_FINISHED, nil)
	var cus, _ = cache.GetCusID(ordRes.CusID)
	s.SendData(ctx, dataRespone(ordRes, emp, cus))
}

func ActionChange(orderID string, empId string, actionStatus common.OrderStatus, extra encode.RawMessage) (ordRes *order.Order) {
	var action = system.NewOrderAction()
	action.EmpId = empId
	fmt.Printf("EMP ID", empId)
	action.OrderID = orderID
	action.Action = actionStatus
	action.Extra = extra
	system.CacheOrderByDay.TriggerTicketAction(action)
	_, err := action.Wait()
	rest.AssertNil(rest.BadRequestValid(err))
	ordRes, _ = system.GetOrderID(orderID)
	return ordRes
}

func (s *OrderServer) handleCancelled(ctx *gin.Context) {
	var cus = oAuth.GetCusFromToken(ctx.Request)
	var body = struct {
		OrderID string `json:"order_id"`
	}{}
	ctx.BindJSON(&body)
	ordRes := ActionChange(body.OrderID, "", common.ORDER_STATUS_CANCELED, nil)
	s.SendData(ctx, dataRespone(ordRes, nil, cus))
}

func (s *OrderServer) handleWorking(ctx *gin.Context) {
	var emp = oAuth.GetEmpFromToken(ctx.Request)
	var body = struct {
		OrderID string `json:"order_id"`
	}{}
	rest.AssertNil(ctx.ShouldBindJSON(&body))
	ordRes := ActionChange(body.OrderID, emp.ID, common.ORDER_STATUS_WORKING, nil)
	var cus, _ = cache.GetCusID(ordRes.CusID)
	s.SendData(ctx, dataRespone(ordRes, emp, cus))
}

func (s *OrderServer) handleRating(ctx *gin.Context) {
	var cus, emp = oAuth.GetUserFromToken(ctx.Request)
	var body *rate.Rate
	ctx.BindJSON(&body)
	var ord, _ = system.GetOrderID(body.OrderID)
	if cus != nil {
		body.CusId = cus.ID
		body.EmpId = ord.EmpID
		body.RateBy = rate.RATE_CUS
		body.CrateRate()
		rest.AssertNil(order.UpdateByRate(body.OrderID))
		ord.Rated = true
		oAuth.UpdateRateToEmp(body.EmpId, body.Rate)
	} else {
		body.CusId = emp.ID
		body.EmpId = ord.EmpID
		body.RateBy = rate.RATE_EMP
		body.CrateRate()
	}
	s.SendData(ctx, nil)
}
