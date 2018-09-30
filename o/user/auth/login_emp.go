package auth

import (
	"ehelp/cache"
	"ehelp/o/user/employee"
	"ehelp/x/rest"
	"errors"
	"net/http"
)

func LoginEmployee(lg *LoginUser) (*employee.Employee, string) {
	validLogin(lg)
	var err, res = employee.GetEmployeeByLogin(lg.Phone, string(lg.Password))
	rest.AssertNil(err)
	return res, CreatePushToken(int(RoleEmployee), res.ID, lg.DeviceId, lg.PushToken).ID
}

func LoginEmployeeFaceBook(lb *LoginFB) (*employee.Employee, string) {
	var err, res = employee.GetEmployeeByLoginFb(lb.FbId)
	rest.AssertNil(err)
	if res != nil {
		var a = FacebookAuth{ID: lb.FbId, Token: lb.FbToken}
		if err := a.IsAuthenticated(); err != nil {
			rest.AssertNil(err)
		}
		rest.AssertNil(res.UpdateEmployeeFb(lb.FbId, lb.FbToken))
		return res, CreatePushToken(int(RoleEmployee), res.ID, lb.DeviceId, lb.PushToken).ID
	}
	return res, ""
}

func LoginEmployeeGmail(lb *LoginGmail) (*employee.Employee, string) {
	var err, res = employee.GetEmployeeByLoginGmail(lb.GmId)
	rest.AssertNil(err)
	if res != nil {
		var a = GmailAuth{GmId: lb.GmId, Token: lb.GmToken}
		if err := a.IsAuthenticated(); err != nil {
			rest.AssertNil(err)
		}
		rest.AssertNil(res.UpdateEmployeeGmail(lb.GmId, lb.GmToken))
		return res, CreatePushToken(int(RoleEmployee), res.ID, lb.DeviceId, lb.PushToken).ID
	}
	return res, ""
}

func CreateEmpFacebook(lb *LoginFB) (*employee.Employee, string) {
	var emp = employee.Employee{}
	emp.Phone = lb.Phone
	emp.FbID = lb.FbId
	emp.FbToken = lb.FbToken
	emp.FullName = lb.FullName
	emp.CrateEmployee()
	return &emp, CreatePushToken(int(RoleEmployee), emp.ID, lb.DeviceId, lb.PushToken).ID
}

func RegisterEmployee(lb *RegisterUser) (*employee.Employee, string) {
	ValidPass(lb.Password)
	var cus = employee.Employee{}
	cus.FullName = lb.FullName
	cus.Phone = lb.Phone
	cus.Password = lb.Password
	var empWork = employee.EmployeeWork{
		AddressWork: lb.Area,
		ServiceIds:  []string{lb.ServiceId},
	}
	cus.EmployeeWork = empWork
	return cus.CrateEmployee(), CreatePushToken(int(RoleEmployee), cus.ID, lb.DeviceId, lb.PushToken).ID
}

func GetEmpFromToken(r *http.Request) *employee.Employee {
	var key = MustGetKey(r)
	if key.Role != int(RoleEmployee) {
		rest.AssertNil(errors.New("Bạn không có quyền này!"))
	}
	var emp, err = cache.GetEmpID(key.UserId)
	rest.AssertNil(err)
	return emp
}

func CreateEmpGmail(lb *LoginGmail) (*employee.Employee, string) {
	var emp = employee.Employee{}
	emp.Phone = lb.Phone
	emp.GmId = lb.GmId
	emp.GmToken = lb.GmToken
	emp.FullName = lb.FullName
	emp.CrateEmployee()
	return &emp, CreatePushToken(int(RoleEmployee), emp.ID, lb.DeviceId, lb.PushToken).ID
}
