package order_server

import (
	//"ehelp/cache"
	//	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/order_hst"
	oAuth "ehelp/o/user/auth"
	"ehelp/x/rest"
	"ehelp/x/web"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (s *OrderServer) handleListOrder(ctx *gin.Context) {
	var query = ctx.Request.URL.Query()
	var status = web.GetArrString("status", ",", query)
	var skip, _ = strconv.ParseInt(query.Get("skip"), 10, 64)
	var limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	var cus, emp = oAuth.GetUserFromToken(ctx.Request)
	var role int
	var addressEmp string
	var services string
	var userId string
	var serviceEmp []string
	if emp != nil {
		addressEmp = emp.EmployeeWork.AddressWork
		role = int(oAuth.RoleEmployee)
		userId = emp.ID
		serviceEmp = emp.EmployeeWork.ServiceIds
	} else {
		role = int(oAuth.RoleCustomer)
		userId = cus.ID
	}
	var data, err = order.GetListOrderByStatus(userId, serviceEmp, addressEmp, role, services, status, int(skip), int(limit))
	rest.AssertNil(rest.WrapBadRequest(err, ""))
	s.SendData(ctx, data)
}

func (s *OrderServer) handleOrderMine(ctx *gin.Context) {
	var emp = oAuth.GetEmpFromToken(ctx.Request)
	var query = ctx.Request.URL.Query()
	var start, _ = strconv.ParseInt(query.Get("start"), 10, 64)
	var end, _ = strconv.ParseInt(query.Get("end"), 10, 64)
	var ords, err = order_hst.GetOrderMine(start, end, emp.ID)
	rest.AssertNil(rest.WrapBadRequest(err, ""))
	var money float32
	for _, item := range ords {
		money += item.MoneyDay
	}
	var data = map[string]interface{}{
		"list_order": ords,
		"employee": map[string]interface{}{
			"rate":          emp.GetMyAvgRate(),
			"all_hour_work": emp.AllHourWork,
			"all_customer":  emp.AllCustomer,
			"rate_5":        emp.Rate5,
		},
		"money_time": money,
	}
	s.SendData(ctx, data)
}
