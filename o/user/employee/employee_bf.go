package employee

import (
	"ehelp/common"
	"ehelp/o/user"
	"ehelp/x/rest"
	"ehelp/x/rest/math"
	"errors"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

var validate = validator.New()

func NewInviteCode() string {
	return strings.TrimPrefix(math.RandStringUpper("", 6), "_")
}

func (emp *Employee) create() error {
	if err := validate.Struct(emp); err != nil {
		return rest.BadRequestValid(err)
	}
	//check  ma gt
	if err := CheckExistByPhone(emp.Phone); err != nil {
		return err
	}
	if strings.Trim(emp.ReferenceCode, " ") != "" {
		var count, _ = CheckEmployeerByReferentCode(emp.ReferenceCode)
		//rest.AssertNil(err)
		if count == 0 {
			return errors.New("Không tồn tại mã giới thiệu!")
		}
	}
	emp.InviteCode = NewInviteCode()
	var psd, err = user.Password(emp.Password).GererateHashedPassword()
	if err != nil {
		return err
	}
	emp.Password = string(psd)
	//emp.BeforeCreate("", 12)
	return nil
}

func CheckEmployeerByReferentCode(refCode string) (int, error) {
	return EmployeeTable.CountWhere(bson.M{
		"invite_code": refCode,
	})
}

func (emp *Employee) delete() {
	emp.BeforeDelete()
}

func (emp *Employee) update() {
	emp.BeforeUpdate()
}

func CheckExistByPhone(phone string) error {
	var count, err = GetEmpByPhone(phone)
	if err != nil && err.Error() != common.NOT_EXIST {
		return err
	}
	if count > 0 {
		return rest.BadRequest("Tài khoản đã tồn tại")
	}
	return nil
}
