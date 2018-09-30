package auth

import (
	"ehelp/o/push_token"
	"ehelp/o/user/customer"
	"ehelp/o/user/employee"
	"ehelp/x/mlog"
	"ehelp/x/rest"
	"errors"
	"g/x/web"
	"net/http"
	"strings"
)

var logAuth = mlog.NewTagLog("Auth")

var LINK_AVATAR string

func MustGetKey(r *http.Request) *push_token.PushToken {
	var token = web.GetToken(r)
	if len(token) < 8 {
		panic(web.Unauthorized("missing or invalid access token"))
	}
	var key, err = push_token.GetByID(token)
	if err != nil {
		panic(rest.Unauthorized("missing or invalid access token"))
	}
	return key
}

func ValidPass(pass string) {
	if len(pass) < 6 {
		rest.AssertNil(rest.BadRequestValid(errors.New("Mật khẩu tối thiểu 6 ký tự")))
	}
}

func validLogin(login *LoginUser) {
	if len(strings.Trim(login.Phone, " ")) == 0 {
		rest.AssertNil(rest.BadRequestValid(errors.New("Nhập tài khoản!")))
	}
	var pass = string(login.Password)
	if len(strings.Trim(pass, " ")) == 0 {
		rest.AssertNil(rest.BadRequestValid(errors.New("Nhập mật khẩu!")))
	}

	ValidPass(string(pass))
}

func DeleteCustomerID(userID string) error {
	return customer.DeleteUserByID(userID)
}

func DeleteEmpID(userID string) error {
	return employee.DeleteUserByID(userID)
}

func UpdateCusNewAndHour(empID string, countCusInOrder int, allHourItemOrder float32) {
	employee.UpdateEmployeeByCusNewAndHour(empID, countCusInOrder, allHourItemOrder)
}

func GetListEmpVsOrderBidding(serviceOrder []string, addressOrder string) ([]string, error) {
	return employee.GetListEmpVsOrderBidding(serviceOrder, addressOrder)
}

func UpdateRateToEmp(empId string, rate int) {
	employee.UpdateRate(empId, rate)
}
