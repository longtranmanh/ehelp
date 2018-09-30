package auth

import (
	"ehelp/cache"
	"ehelp/common"
	"ehelp/o/push_token"
	"ehelp/o/user/customer"
	"ehelp/o/user/employee"
	"ehelp/x/rest"
	"errors"
	"net/http"
)

type LoginUser struct {
	Phone     string   `json:"phone"`
	Password  Password `json:"password"`
	DeviceId  string   `json:"device_id"`
	PushToken string   `json:"push_token"`
}

type LoginFB struct {
	FbId      string `json:"fb_id"`
	FbToken   string `json:"fb_token"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	DeviceId  string `json:"device_id"`
	PushToken string `json:"push_token"`
}

type LoginGmail struct {
	GmId      string `json:"gm_id"`
	GmToken   string `json:"gm_token"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	DeviceId  string `json:"device_id"`
	PushToken string `json:"push_token"`
}

type RegisterUser struct {
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
	ReferenCode string `json:"reference_code"`
	Area        string `json:"area"`
	ServiceId   string `json:"service_id"`
	DeviceId    string `json:"device_id"`
	PushToken   string `json:"push_token"`
}

type Password string
type Token string

func CreatePushToken(role int, userId string, deviceID string, pushToken string) *push_token.PushToken {
	var psh = push_token.PushToken{
		Role:      role,
		UserId:    userId,
		DeviceId:  deviceID,
		PushToken: pushToken,
	}
	return psh.CratePushToken()
}

func LoginCustomer(lg *LoginUser) (*customer.Customer, string) {
	validLogin(lg)
	var err, res = customer.GetCustomerByLogin(lg.Phone, string(lg.Password))
	rest.AssertNil(err)
	return res, CreatePushToken(int(RoleCustomer), res.ID, lg.DeviceId, lg.PushToken).ID
}

func LoginCustomerFaceBook(lb *LoginFB) (*customer.Customer, string) {
	var err, res = customer.GetCustomerByLoginFb(lb.FbId)
	if err != nil {
		if err.Error() == common.NOT_EXIST {
			rest.AssertNil(rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!")))
		}
		rest.AssertNil(err)
	}

	if res != nil {
		var a = FacebookAuth{ID: lb.FbId, Token: lb.FbToken}
		if err := a.IsAuthenticated(); err != nil {
			rest.AssertNil(err)
		}
		rest.AssertNil(res.UpdateCustomerFb(lb.FbId, lb.FbToken))
		return res, CreatePushToken(int(RoleCustomer), res.ID, lb.DeviceId, lb.PushToken).ID
	}
	return res, ""
}

func LoginCustomerGmail(lb *LoginGmail) (*customer.Customer, string) {
	var err, res = customer.GetCustomerByLoginGmail(lb.GmId)
	rest.AssertNil(err)
	if res != nil {
		var a = GmailAuth{GmId: lb.GmId, Token: lb.GmToken}
		if err := a.IsAuthenticated(); err != nil {
			rest.AssertNil(err)
		}
		rest.AssertNil(res.UpdateCustomerGmail(lb.GmId, lb.GmToken))
		return res, CreatePushToken(int(RoleCustomer), res.ID, lb.DeviceId, lb.PushToken).ID
	}
	return res, ""
}

func RegisterCustomer(lb *RegisterUser) (*customer.Customer, string) {
	ValidPass(lb.Password)
	var cus = customer.Customer{}
	cus.FullName = lb.FullName
	cus.Phone = lb.Phone
	cus.Email = lb.Email
	cus.Password = lb.Password
	return cus.CrateCustomer(), CreatePushToken(int(RoleCustomer), cus.ID, lb.DeviceId, lb.PushToken).ID
}

func CreateCusFacebook(lb *LoginFB) (*customer.Customer, string) {
	var cus = customer.Customer{}
	cus.Phone = lb.Phone
	cus.FbID = lb.FbId
	cus.FbToken = lb.FbToken
	cus.FullName = lb.FullName
	cus.CrateCustomer()
	return &cus, CreatePushToken(int(RoleCustomer), cus.ID, lb.DeviceId, lb.PushToken).ID
}

func CreateCusGmail(lb *LoginGmail) (*customer.Customer, string) {
	var cus = customer.Customer{}
	cus.Phone = lb.Phone
	cus.GmId = lb.GmId
	cus.GmToken = lb.GmToken
	cus.FullName = lb.FullName
	cus.CrateCustomer()
	return &cus, CreatePushToken(int(RoleCustomer), cus.ID, lb.DeviceId, lb.PushToken).ID
}

func Logout(token string) error {
	return push_token.UpdatePushToken(token)
}

func GetCusFromToken(r *http.Request) *customer.Customer {
	var key = MustGetKey(r)
	if key.Role != int(RoleCustomer) {
		rest.AssertNil(errors.New("Bạn không có quyền này!"))
	}
	var cus, err = cache.GetCusID(key.UserId)
	rest.AssertNil(err)
	return cus
}

func GetUserFromToken(r *http.Request) (*customer.Customer, *employee.Employee) {
	var key = MustGetKey(r)
	var emp, err = cache.GetEmpID(key.UserId)
	if emp == nil {
		var cus, err = cache.GetCusID(key.UserId)
		rest.AssertNil(err)
		return cus, nil
	}
	rest.AssertNil(err)
	return nil, emp
}
